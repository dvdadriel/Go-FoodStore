package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-food-store/helpers"
)

// env mengembalikan nilai environment variable, atau fallback kalau tidak
// di-set. String kosong yang di-set secara eksplisit dianggap nilai yang
// disengaja — password kosong itu sah untuk MySQL root di lokal.
func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

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
