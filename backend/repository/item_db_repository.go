package repository

import (
	"errors"
	"gin-market/models"
	"gorm.io/gorm"
	"log"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindMyAll(userId uint) (*[]models.Item, error)
	FindById(targetId uint, userId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	Delete(targetId uint, userId uint) error
}

type ItemDBRepository struct {
	db *gorm.DB
}

func NewItemDBRepository(db *gorm.DB) IItemRepository {
	return &ItemDBRepository{
		db: db,
	}
}

func (ir *ItemDBRepository) FindAll() (*[]models.Item, error) {
	items := make([]models.Item, 0)

	result := ir.db.Find(&items)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}

	log.Println(items)
	return &items, nil
}

func (ir *ItemDBRepository) FindMyAll(userId uint) (*[]models.Item, error) {
	items := make([]models.Item, 0)

	result := ir.db.Find(&items, "user_id = ?", userId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}

	return &items, nil
}

func (ir *ItemDBRepository) FindById(targetId uint, userId uint) (*models.Item, error) {
	var item models.Item

	result := ir.db.First(&item, "id = ? AND user_id = ?", targetId, userId)
	if result.Error != nil {
		log.Println(result.Error.Error())

		if result.Error.Error() == "record not found" {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}

	log.Printf("successfully called FindById = %v", item)
	return &item, nil
}

func (ir *ItemDBRepository) Create(newItem models.Item) (*models.Item, error) {
	result := ir.db.Create(&newItem)
	if result.Error != nil {
		log.Printf("failed to Create = %s", result.Error.Error())
		return nil, result.Error
	}

	log.Printf("successfully called Create = %v", newItem)
	return &newItem, nil
}

func (ir *ItemDBRepository) Update(updateItem models.Item) (*models.Item, error) {
	result := ir.db.Save(&updateItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateItem, nil
}

func (ir *ItemDBRepository) Delete(targetId uint, userId uint) error {
	targetItem, err := ir.FindById(targetId, userId)
	if err != nil {
		return err
	}

	// [if delete completely]
	//result := ir.db.Unscoped().Delete(&targetItem)
	result := ir.db.Delete(&targetItem)
	if result.Error != nil {
		log.Printf("failed to delete = %v", targetItem)
		return result.Error
	}

	log.Printf("successfully Deleted = %v", targetItem)
	return nil
}
