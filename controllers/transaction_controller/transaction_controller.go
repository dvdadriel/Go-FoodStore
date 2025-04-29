package transactioncontroller

import "net/http"

type TransactionController interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request)
	UpdateTransaction(w http.ResponseWriter, r *http.Request)
	DeleteTransaction(w http.ResponseWriter, r *http.Request)
	GetTransactionById(w http.ResponseWriter, r *http.Request)
	FindAllPaidTransaction(w http.ResponseWriter, r *http.Request)
	FindAllUnpaidTransaction(w http.ResponseWriter, r *http.Request)
	AcceptTransactionPayment(w http.ResponseWriter, r *http.Request)
}
