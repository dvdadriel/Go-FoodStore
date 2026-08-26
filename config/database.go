package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-food-store/helpers"
)

// DSN membentuk DSN MySQL dari environment, dengan default yang cocok untuk
// pengembangan lokal.
func DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env("DB_USER", "root"),
		env("DB_PASSWORD", ""),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "go-food-store"),
	)
}

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open(DSN()))
	helpers.PanicHelper(err)
	return db
}
