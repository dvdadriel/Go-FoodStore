package transactioncontroller

import (
	"go-food-store/exception"
	"go-food-store/helpers"
	"go-food-store/json/request"
	transactionservice "go-food-store/services/transaction_service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TransactionControllerImpl struct {
	TransactionService transactionservice.TransactionService
}

func NewTransactionController(TransactionService transactionservice.TransactionService) TransactionController {
	return &TransactionControllerImpl{
		TransactionService: TransactionService,
	}
}

// AcceptTransactionPayment implements TransactionController.
func (t *TransactionControllerImpl) AcceptTransactionPayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["transactionId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := t.TransactionService.AcceptPayment(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// CreateTransaction implements TransactionController.
func (t *TransactionControllerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	newTransaction := &request.CreateTransactionReq{}
	err := helpers.Ummarshal(r, newTransaction)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := t.TransactionService.Create(*newTransaction)
		helpers.WriteJSON(w, response)
	}
}

// DeleteTransaction implements TransactionController.
func (t *TransactionControllerImpl) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["transactionId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := t.TransactionService.Delete(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// FindAllPaidTransaction implements TransactionController.
func (t *TransactionControllerImpl) FindAllPaidTransaction(w http.ResponseWriter, r *http.Request) {
	response := t.TransactionService.FindAllPaid()
	helpers.WriteJSON(w, response)
}

// FindAllUnpaidTransaction implements TransactionController.
func (t *TransactionControllerImpl) FindAllUnpaidTransaction(w http.ResponseWriter, r *http.Request) {
	response := t.TransactionService.FindAllUnpaid()
	helpers.WriteJSON(w, response)
}

// GetTransactionById implements TransactionController.
func (t *TransactionControllerImpl) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["transactionId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := t.TransactionService.FindById(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// UpdateTransaction implements TransactionController.
func (t *TransactionControllerImpl) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	updatedTransaction := &request.UpdateTransactionReq{}
	err := helpers.Ummarshal(r, updatedTransaction)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		params := mux.Vars(r)
		id, err := strconv.ParseInt(params["transactionId"], 0, 0)
		if err != nil {
			exception.BadRequestErr(w, r)
		} else {
			updatedTransaction.Id = uint(id)
			response := t.TransactionService.Update(*updatedTransaction)
			helpers.WriteJSON(w, response)
		}
	}
}
