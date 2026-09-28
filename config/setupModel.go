package config

import (
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"go-food-store/auth"
	authcontroller "go-food-store/controllers/auth_controller"
	customercontroller "go-food-store/controllers/customer_controller"
	foodcontroller "go-food-store/controllers/food_controller"
	transactioncontroller "go-food-store/controllers/transaction_controller"
	"go-food-store/middleware"
	customerrepository "go-food-store/repositories/customer_repository"
	foodrepository "go-food-store/repositories/food_repository"
	transactionrepository "go-food-store/repositories/transaction_repository"
	userrepository "go-food-store/repositories/user_repository"
	"go-food-store/routes"
	authservice "go-food-store/services/auth_service"
	customerservice "go-food-store/services/customer_service"
	foodservice "go-food-store/services/food_service"
	transactionservice "go-food-store/services/transaction_service"
)

func SetupModel(db *gorm.DB, validate *validator.Validate) *mux.Router {
	ttl := TokenTTL()
	issuer := auth.NewIssuer(JWTSecret(), ttl)

	userRepo := userrepository.NewUserRepo(db)
	authService := authservice.NewAuthService(userRepo, issuer, int64(ttl.Seconds()), validate)
	authController := authcontroller.NewAuthController(authService)

	custRepo := customerrepository.NewCustomerRepo(db)
	custService := customerservice.NewCustService(custRepo, validate)
	custController := customercontroller.NewCustomerController(custService)

	foodRepo := foodrepository.NewFoodRepo(db)
	foodService := foodservice.NewFoodService(foodRepo, validate)
	foodController := foodcontroller.NewFoodController(foodService)

	transactionRepo := transactionrepository.NewTransactionRepository(db)
	transactionService := transactionservice.NewTransactionService(transactionRepo, validate)
	transactionController := transactioncontroller.NewTransactionController(transactionService)

	router := mux.NewRouter()

	// Urutan berarti: Recover terluar supaya ia juga menangkap panic dari
	// middleware di dalamnya, dan Log tepat di dalamnya supaya request yang
	// berakhir panic tetap tercatat dengan status 500-nya.
	router.Use(middleware.Recover)
	router.Use(middleware.Log)
	router.Use(middleware.CORS(CORSOrigin()))

	routes.NewRouter(router, db, issuer, authController, custController, foodController, transactionController)
	return router
}
