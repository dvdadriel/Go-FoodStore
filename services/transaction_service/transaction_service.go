package transactionservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
)

type TransactionService interface {
	Create(transaction request.CreateTransactionReq) response.WebResponse
	Update(transaction request.UpdateTransactionReq) response.WebResponse
	Delete(transactionId uint) response.WebResponse
	FindAllPaid(page request.Page) response.WebResponse
	FindAllUnpaid(page request.Page) response.WebResponse
	AcceptPayment(transactionId uint) response.WebResponse
	FindById(transactionId uint) response.WebResponse
}
