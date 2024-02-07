package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string `gorm:"notnull;unique"`
	Password string `gorm:"notnull"`
	items    []Item `gorm:"constraint:OnDelete:CASCADE"`
}
