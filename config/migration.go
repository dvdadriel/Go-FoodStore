package config

import (
	"go-food-store/helpers"
	"go-food-store/models"

	"gorm.io/gorm"
)

func MigrateAllTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{},
		&models.Food{},
		&models.Customer{},
		&models.Transaction{},
		&models.Transaction_Food{})
	helpers.PanicHelper(err)
}
