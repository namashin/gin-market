package repository

import (
	"errors"
	"gin-market/models"
)

type ItemMemoryRepository struct {
	items []models.Item
}

func NewItemMemoryRepository(items []models.Item) IItemRepository {
	return &ItemMemoryRepository{
		items: items,
	}
}

func (imr *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
	return &imr.items, nil
}

func (imr *ItemMemoryRepository) FindById(targetId uint, userId uint) (*models.Item, error) {
	for _, item := range imr.items {
		if item.ID == targetId && item.UserID == userId {
			return &item, nil
		}
	}

	return nil, errors.New("item not found")
}

func (imr *ItemMemoryRepository) Create(newItem models.Item) (*models.Item, error) {
	newItem.ID = uint(len(imr.items)) + 1
	imr.items = append(imr.items, newItem)
	return &newItem, nil
}

func (imr *ItemMemoryRepository) Update(updateItem models.Item) (*models.Item, error) {
	for i, item := range imr.items {
		if item.ID == updateItem.ID {
			imr.items[i] = updateItem
			return &imr.items[i], nil
		}
	}

	return nil, errors.New("target item not found")
}

func (imr *ItemMemoryRepository) Delete(targetId uint, userId uint) error {
	for i, item := range imr.items {
		if item.ID == targetId && item.UserID == userId {
			// remove i-index item
			imr.items = append(imr.items[:i], imr.items[i+1:]...)
			return nil
		}
	}

	return errors.New("item not found")
}

func (imr *ItemMemoryRepository) FindMyAll(userId uint) (*[]models.Item, error) {
	items := make([]models.Item, 0)

	for _, item := range imr.items {
		if item.UserID == userId {
			items = append(items, item)
		}
	}

	return &items, nil
}
