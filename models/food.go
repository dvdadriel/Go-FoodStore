package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	FoodName    string        `gorm:"not null" json:"FoodName"`
	FoodPrice   float64       `gorm:"not null" json:"FoodPrice"`
	Transaction []Transaction `gorm:"many2many:transaction_foods;" json:"Transaction"`
}
