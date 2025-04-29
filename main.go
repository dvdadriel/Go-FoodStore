package main

import (
	"go-food-store/config"
	"go-food-store/helpers"
	"net/http"

	"github.com/go-playground/validator"
)

func main() {
	db := config.ConnectToDatabase()
	validate := validator.New()

	config.MigrateAllTable(db)

	router := config.SetupModel(db, validate)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	helpers.PanicHelper(err)
}
