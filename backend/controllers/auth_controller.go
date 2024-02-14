package controllers

import (
	"gin-market/dto"
	"gin-market/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IAuthController interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) IAuthController {
	return &AuthController{
		service: service,
	}
}

func (ac *AuthController) SignUp(ctx *gin.Context) {
	var sui dto.SignUpInput

	err := ctx.ShouldBindJSON(&sui)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.service.SignUp(sui.Email, sui.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": sui})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var li dto.LoginInput

	err := ctx.ShouldBindJSON(&li)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.service.Login(li.Email, li.Password)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
