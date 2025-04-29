package config

import (
	customercontroller "go-food-store/controllers/customer_controller"
	foodcontroller "go-food-store/controllers/food_controller"
	transactioncontroller "go-food-store/controllers/transaction_controller"
	customerrepository "go-food-store/repositories/customer_repository"
	foodrepository "go-food-store/repositories/food_repository"
	transactionrepository "go-food-store/repositories/transaction_repository"
	"go-food-store/routes"
	customerservice "go-food-store/services/customer_service"
	foodservice "go-food-store/services/food_service"
	transactionservice "go-food-store/services/transaction_service"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupModel(db *gorm.DB, validate *validator.Validate) *mux.Router {
	custRepo := customerrepository.NewCustomerRepo(db)
	custService := customerservice.NewCustService(custRepo, validate)
	custController := customercontroller.NewCustomerController(custService)

	foodRepo := foodrepository.NewFoodRepo(db)
	foodService := foodservice.NewFoodService(foodRepo, validate)
	foodcontroller := foodcontroller.NewFoodController(foodService)

	transactionRepo := transactionrepository.NewTransactionRepository(db)
	transactionService := transactionservice.NewTransactionService(transactionRepo, validate)
	transactionController := transactioncontroller.NewTransactionController(transactionService)

	router := mux.NewRouter()

	routes.NewRouter(router, custController, foodcontroller, transactionController)
	return router
}
