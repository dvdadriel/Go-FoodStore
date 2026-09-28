package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"go-food-store/auth"
	authcontroller "go-food-store/controllers/auth_controller"
	customercontroller "go-food-store/controllers/customer_controller"
	foodcontroller "go-food-store/controllers/food_controller"
	transactioncontroller "go-food-store/controllers/transaction_controller"
	"go-food-store/exception"
	"go-food-store/middleware"
	"go-food-store/models"
)

// NewRouter memetakan seluruh endpoint sekaligus menentukan siapa yang boleh
// memanggilnya.
//
// Dibagi tiga kelompok yang urutan pendaftarannya berarti: mux mencoba route
// sesuai urutan, jadi route publik harus terdaftar sebelum subrouter yang
// berpenjaga, kalau tidak penjaganya akan menangkap request yang mestinya
// terbuka.
func NewRouter(
	router *mux.Router,
	db *gorm.DB,
	issuer *auth.Issuer,
	authCtrl authcontroller.AuthController,
	cust customercontroller.CustomerController,
	food foodcontroller.FoodController,
	transaction transactioncontroller.TransactionController,
) {
	// --- Terbuka untuk umum ---
	router.HandleFunc("/health", Health(db)).Methods(http.MethodGet)
	router.HandleFunc("/auth/register", authCtrl.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authCtrl.Login).Methods(http.MethodPost)

	// Daftar menu memang untuk dilihat siapa saja — itu etalase toko.
	router.HandleFunc("/food/", food.GetAllFood).Methods(http.MethodGet)
	router.HandleFunc("/food/{foodId}", food.GetFoodById).Methods(http.MethodGet)

	// --- Hanya admin: mengelola toko ---
	admin := router.NewRoute().Subrouter()
	admin.Use(middleware.RequireAuth(issuer), middleware.RequireRole(models.RoleAdmin))

	admin.HandleFunc("/food/", food.CreateFood).Methods(http.MethodPost)
	admin.HandleFunc("/food/{foodId}", food.UpdateFood).Methods(http.MethodPut)
	admin.HandleFunc("/food/{foodId}", food.DeleteFood).Methods(http.MethodDelete)

	admin.HandleFunc("/cust/", cust.GetAllCust).Methods(http.MethodGet)
	admin.HandleFunc("/cust/", cust.CreateCust).Methods(http.MethodPost)
	admin.HandleFunc("/cust/{custId}", cust.GetCustById).Methods(http.MethodGet)
	admin.HandleFunc("/cust/{custId}", cust.UpdateCust).Methods(http.MethodPut)
	admin.HandleFunc("/cust/{custId}", cust.DeleteCust).Methods(http.MethodDelete)

	// Pembukuan dan penerimaan pembayaran adalah pekerjaan kasir.
	admin.HandleFunc("/transaction/paid/", transaction.FindAllPaidTransaction).Methods(http.MethodGet)
	admin.HandleFunc("/transaction/unpaid/", transaction.FindAllUnpaidTransaction).Methods(http.MethodGet)
	admin.HandleFunc("/transaction/acc/{transactionId}", transaction.AcceptTransactionPayment).Methods(http.MethodPut)

	// --- Perlu login, peran apa pun ---
	authenticated := router.NewRoute().Subrouter()
	authenticated.Use(middleware.RequireAuth(issuer))

	authenticated.HandleFunc("/transaction/", transaction.CreateTransaction).Methods(http.MethodPost)
	authenticated.HandleFunc("/transaction/{transactionId}", transaction.GetTransactionById).Methods(http.MethodGet)
	authenticated.HandleFunc("/transaction/{transactionId}", transaction.UpdateTransaction).Methods(http.MethodPut)
	authenticated.HandleFunc("/transaction/{transactionId}", transaction.DeleteTransaction).Methods(http.MethodDelete)

	router.MethodNotAllowedHandler = http.HandlerFunc(exception.NotAllowedErr)
	router.NotFoundHandler = http.HandlerFunc(exception.NotFoundErr)
}
