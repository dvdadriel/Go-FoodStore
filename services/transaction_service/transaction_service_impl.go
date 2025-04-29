package transactionservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	transaction "go-food-store/repositories/transaction_repository"
	"net/http"

	"github.com/go-playground/validator"
)

type TransactionServiceImpl struct {
	TransactionRepo transaction.TransactionRepositories
	Validate        *validator.Validate
}

func NewTransactionService(TransactionRepo transaction.TransactionRepositories, Validate *validator.Validate) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepo: TransactionRepo,
		Validate:        Validate,
	}
}

// AcceptPayment implements TransactionService.
func (t *TransactionServiceImpl) AcceptPayment(transactionId uint) response.WebResponse {
	response := t.TransactionRepo.AcceptPayment(transactionId)
	return response
}

// Create implements TransactionService.
func (t *TransactionServiceImpl) Create(transaction request.CreateTransactionReq) response.WebResponse {
	err := t.Validate.Struct(transaction)
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad request",
			Message: "Please check the request",
			Data:    nil,
		}
	}
	response, found := t.TransactionRepo.GetCustomerById(transaction.CustomerId)
	if !found {
		return response
	}
	for _, detail := range transaction.Food {
		response, found := t.TransactionRepo.GetFoodById(detail.FoodId)
		if !found {
			return response
		}
	}
	response = t.TransactionRepo.CreateTransaction(transaction)
	return response
}

// Delete implements TransactionService.
func (t *TransactionServiceImpl) Delete(transactionId uint) response.WebResponse {
	response, found := t.TransactionRepo.GetTransactionById(transactionId)
	if !found {
		return response
	}
	response = t.TransactionRepo.DeleteTransaction(transactionId)
	return response
}

// FindAllPaid implements TransactionService.
func (t *TransactionServiceImpl) FindAllPaid() response.WebResponse {
	response := t.TransactionRepo.GetAllPaidTransaction()
	return response
}

// FindAllUnpaid implements TransactionService.
func (t *TransactionServiceImpl) FindAllUnpaid() response.WebResponse {
	response := t.TransactionRepo.GetUnpaidTransaction()
	return response
}

// FindById implements TransactionService.
func (t *TransactionServiceImpl) FindById(transactionId uint) response.WebResponse {
	response, _ := t.TransactionRepo.GetTransactionById(transactionId)
	return response
}

// Update implements TransactionService.
func (t *TransactionServiceImpl) Update(transaction request.UpdateTransactionReq) response.WebResponse {
	err := t.Validate.Struct(transaction)
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad request",
			Message: "Please check the request",
			Data:    nil,
		}
	}

	response := t.TransactionRepo.UpdateTransaction(transaction)
	return response
}
