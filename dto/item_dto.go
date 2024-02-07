package dto

// Data Transfer Object

type CreateItemInput struct {
	Name        string `json:"name" binding:"required,min=2"`
	Price       uint   `json:"price" binding:"required,min=1,max=99999"`
	Description string `json:"description"`
}

type UpdateItemInput struct {
	Name        *string `json:"name" binding:"omitempty,min=2"`
	Price       *uint   `json:"price" binding:"omitempty,min=1,max=99999"`
	Description *string `json:"description"`
	SoldOut     *bool   `json:"soldOut"`
}
