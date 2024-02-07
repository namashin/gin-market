package controllers

import (
	"gin-market/dto"
	"gin-market/models"
	"gin-market/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IItemController interface {
	FindAll(ctx *gin.Context)
	FindMyAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type ItemController struct {
	service services.IItemService
}

func NewItemController(service services.IItemService) IItemController {
	return &ItemController{
		service: service,
	}
}

func (ic *ItemController) FindAll(ctx *gin.Context) {
	items, err := ic.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (ic *ItemController) FindMyAll(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := user.(*models.User).ID

	items, err := ic.service.FindMyAll(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusFound, gin.H{"data": items})
}

func (ic *ItemController) FindById(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := user.(*models.User).ID

	itemUint64Id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetItem, err := ic.service.FindById(uint(itemUint64Id), userId)
	if err != nil {
		if err.Error() == "item not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusFound, gin.H{"data": targetItem})
}

func (ic *ItemController) Create(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := user.(*models.User).ID

	var cii dto.CreateItemInput

	err := ctx.ShouldBindJSON(&cii)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := ic.service.Create(cii, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newItem})
}

func (ic *ItemController) Update(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := user.(*models.User).ID

	itemUint64Id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var uii dto.UpdateItemInput

	// bind uii with json
	err = ctx.ShouldBindJSON(&uii)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedItem, err := ic.service.Update(uint(itemUint64Id), userId, uii)
	if err != nil {
		if err.Error() == "item not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedItem})
}

func (ic *ItemController) Delete(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := user.(*models.User).ID

	itemUint64Id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ic.service.Delete(uint(itemUint64Id), userId)
	if err != nil {
		if err.Error() == "item not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
