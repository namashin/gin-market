package services

import (
	"gin-market/dto"
	"gin-market/models"
	"gin-market/repository"
)

type IItemService interface {
	FindAll() (*[]models.Item, error)
	FindMyAll(userId uint) (*[]models.Item, error)
	FindById(targetId uint, userId uint) (*models.Item, error)
	Create(createItemInput dto.CreateItemInput, userId uint) (*models.Item, error)
	Update(targetId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error)
	Delete(targetId uint, userId uint) error
}

type ItemService struct {
	repository repository.IItemRepository
}

func NewItemService(repository repository.IItemRepository) IItemService {
	return &ItemService{repository: repository}
}

func (is *ItemService) FindAll() (*[]models.Item, error) {
	return is.repository.FindAll()
}

func (is *ItemService) FindMyAll(userId uint) (*[]models.Item, error) {
	return is.repository.FindMyAll(userId)
}

func (is *ItemService) FindById(targetId uint, userId uint) (*models.Item, error) {
	return is.repository.FindById(targetId, userId)
}

func (is *ItemService) Create(createItemInput dto.CreateItemInput, userId uint) (*models.Item, error) {
	newItem := models.Item{
		Name:        createItemInput.Name,
		Price:       createItemInput.Price,
		Description: createItemInput.Description,
		SoldOut:     false,
		UserID:      userId,
	}

	return is.repository.Create(newItem)
}

func (is *ItemService) Update(targetId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error) {
	// find update item
	targetItem, err := is.FindById(targetId, userId)
	if err != nil {
		return nil, err
	}

	// update value if not nil
	if updateItemInput.Name != nil {
		targetItem.Name = *updateItemInput.Name
	}
	if updateItemInput.Price != nil {
		targetItem.Price = *updateItemInput.Price
	}
	if updateItemInput.Description != nil {
		targetItem.Description = *updateItemInput.Description
	}
	if updateItemInput.SoldOut != nil {
		targetItem.SoldOut = *updateItemInput.SoldOut
	}

	return is.repository.Update(*targetItem)
}

func (is *ItemService) Delete(targetId uint, userId uint) error {
	return is.repository.Delete(targetId, userId)
}
