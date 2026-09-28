package transactionservice

import (
	"net/http"

	"github.com/go-playground/validator"

	"go-food-store/json/request"
	"go-food-store/json/response"
	transaction "go-food-store/repositories/transaction_repository"
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

func badRequest(message string) response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusBadRequest,
		Status:  "Bad Request",
		Message: message,
		Data:    nil,
	}
}

// duplicateFood melaporkan apakah ada FoodId yang muncul lebih dari sekali.
//
// Baris detail memakai (transaction_id, food_id) sebagai primary key, jadi
// dua baris untuk makanan yang sama menabrak constraint dan dulu keluar
// sebagai 500. Itu kesalahan client, bukan kesalahan server: gabungkan
// quantity-nya di sisi pemanggil.
func duplicateFood(items []request.FoodDetailReq) bool {
	seen := make(map[uint]struct{}, len(items))
	for _, item := range items {
		if _, exists := seen[item.FoodId]; exists {
			return true
		}
		seen[item.FoodId] = struct{}{}
	}
	return false
}

// AcceptPayment implements TransactionService.
func (t *TransactionServiceImpl) AcceptPayment(transactionId uint) response.WebResponse {
	return t.TransactionRepo.AcceptPayment(transactionId)
}

// Create implements TransactionService.
func (t *TransactionServiceImpl) Create(req request.CreateTransactionReq) response.WebResponse {
	if err := t.Validate.Struct(req); err != nil {
		return badRequest("Please check the request")
	}
	if duplicateFood(req.Food) {
		return badRequest("Duplicate FoodId in one transaction, merge the quantity instead")
	}
	// Keberadaan makanan diperiksa di dalam transaksi database oleh
	// repository, jadi tidak diulang di sini: cek ganda hanya menambah
	// perjalanan ke database dan tetap bisa basi sebelum insert berjalan.
	if res, found := t.TransactionRepo.GetCustomerById(req.CustomerId); !found {
		return res
	}
	return t.TransactionRepo.CreateTransaction(req)
}

// Delete implements TransactionService.
func (t *TransactionServiceImpl) Delete(transactionId uint) response.WebResponse {
	return t.TransactionRepo.DeleteTransaction(transactionId)
}

// FindAllPaid implements TransactionService.
func (t *TransactionServiceImpl) FindAllPaid(page request.Page) response.WebResponse {
	return t.TransactionRepo.GetAllPaidTransaction(page)
}

// FindAllUnpaid implements TransactionService.
func (t *TransactionServiceImpl) FindAllUnpaid(page request.Page) response.WebResponse {
	return t.TransactionRepo.GetUnpaidTransaction(page)
}

// FindById implements TransactionService.
func (t *TransactionServiceImpl) FindById(transactionId uint) response.WebResponse {
	res, _ := t.TransactionRepo.GetTransactionById(transactionId)
	return res
}

// Update implements TransactionService.
func (t *TransactionServiceImpl) Update(req request.UpdateTransactionReq) response.WebResponse {
	if err := t.Validate.Struct(req); err != nil {
		return badRequest("Please check the request")
	}
	if duplicateFood(req.Food) {
		return badRequest("Duplicate FoodId in one transaction, merge the quantity instead")
	}
	if res, found := t.TransactionRepo.GetCustomerById(req.CustomerId); !found {
		return res
	}
	return t.TransactionRepo.UpdateTransaction(req)
}
