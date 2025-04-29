package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	CustomerId uint     `json:"CustomerId"`
	Customer   Customer `gorm:"foreignKey:CustomerId" json:"Customer"`
	Food       []Food   `gorm:"many2many:transaction_foods;" json:"Food"`
	AlreadyPay bool     `gorm:"default:false" json:"AlreadyPay"`
}
