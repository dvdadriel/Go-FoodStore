package config

import (
	"go-food-store/helpers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go-food-store?charset=utf8mb4&parseTime=True&loc=Local"))
	helpers.PanicHelper(err)
	return db
}
