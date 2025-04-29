package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	CustomerName string `gorm:"not null" json:"CustomerName"`
}
