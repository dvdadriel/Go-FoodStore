package routes

import (
	customercontroller "go-food-store/controllers/customer_controller"
	foodcontroller "go-food-store/controllers/food_controller"
	transactioncontroller "go-food-store/controllers/transaction_controller"
	"go-food-store/exception"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(router *mux.Router, db *gorm.DB, cust customercontroller.CustomerController, food foodcontroller.FoodController, transaction transactioncontroller.TransactionController) {
	router.HandleFunc("/health", Health(db)).Methods("GET")

	router.HandleFunc("/cust/", cust.GetAllCust).Methods("GET")
	router.HandleFunc("/cust/", cust.CreateCust).Methods("POST")
	router.HandleFunc("/cust/{custId}", cust.GetCustById).Methods("GET")
	router.HandleFunc("/cust/{custId}", cust.UpdateCust).Methods("PUT")
	router.HandleFunc("/cust/{custId}", cust.DeleteCust).Methods("DELETE")

	router.HandleFunc("/food/", food.GetAllFood).Methods("GET")
	router.HandleFunc("/food/", food.CreateFood).Methods("POST")
	router.HandleFunc("/food/{foodId}", food.GetFoodById).Methods("GET")
	router.HandleFunc("/food/{foodId}", food.UpdateFood).Methods("PUT")
	router.HandleFunc("/food/{foodId}", food.DeleteFood).Methods("DELETE")

	router.HandleFunc("/transaction/paid/", transaction.FindAllPaidTransaction).Methods("GET")
	router.HandleFunc("/transaction/unpaid/", transaction.FindAllUnpaidTransaction).Methods("GET")
	router.HandleFunc("/transaction/", transaction.CreateTransaction).Methods("POST")
	router.HandleFunc("/transaction/{transactionId}", transaction.GetTransactionById).Methods("GET")
	router.HandleFunc("/transaction/{transactionId}", transaction.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transaction/{transactionId}", transaction.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/transaction/acc/{transactionId}", transaction.AcceptTransactionPayment).Methods("PUT")

	router.MethodNotAllowedHandler = http.HandlerFunc(exception.NotAllowedErr)
	router.NotFoundHandler = http.HandlerFunc(exception.NotFoundErr)
}
