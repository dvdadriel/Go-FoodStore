package transaction

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
)

type TransactionRepositories interface {
	CreateTransaction(transaction request.CreateTransactionReq) response.WebResponse
	UpdateTransaction(transaction request.UpdateTransactionReq) response.WebResponse
	DeleteTransaction(TransactionId uint) response.WebResponse
	GetAllPaidTransaction() response.WebResponse
	GetUnpaidTransaction() response.WebResponse
	GetTransactionById(TransactionId uint) (response.WebResponse, bool)
	AcceptPayment(TransactionId uint) response.WebResponse
	GetCustomerById(CustomerId uint) (response.WebResponse, bool)
	GetFoodById(FoodId uint) (response.WebResponse, bool)
}
