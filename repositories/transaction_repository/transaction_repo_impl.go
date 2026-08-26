package transaction

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type TransactionRepositoriesImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepositories {
	return &TransactionRepositoriesImpl{
		DB: DB,
	}
}

// GetCustomerById implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetCustomerById(CustomerId uint) (response.WebResponse, bool) {
	var customer models.Customer
	err := t.DB.Where("id = ?", CustomerId).Find(&customer).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get customer due to server error",
			Data:    nil,
		}, false
	} else if customer.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Customer with id " + strconv.FormatUint(uint64(CustomerId), 10) + " not found",
		}, false
	}
	return response.WebResponse{}, true
}

// GetFodById implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetFoodById(FoodId uint) (response.WebResponse, bool) {
	var food models.Food
	err := t.DB.Where("id = ?", FoodId).Find(&food).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get food due to server error",
			Data:    nil,
		}, false
	} else if food.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Food with id " + strconv.FormatUint(uint64(FoodId), 10) + " not found",
			Data:    nil,
		}, false
	}
	return response.WebResponse{}, true
}

// AcceptPayment implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) AcceptPayment(TransactionId uint) response.WebResponse {
	var transaction models.Transaction
	err := t.DB.Where("id = ?", TransactionId).Find(&transaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction due to server error",
			Data:    nil,
		}
	} else if transaction.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Transaction not found",
			Data:    nil,
		}
	}
	transaction.AlreadyPay = true
	err = t.DB.Updates(transaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update transaction status due to server error",
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully update transaction status",
		Data:    nil,
	}
}

// CreateTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) CreateTransaction(transaction request.CreateTransactionReq) response.WebResponse {
	newTransaction := models.Transaction{
		CustomerId: transaction.CustomerId,
		AlreadyPay: false,
	}
	err := t.DB.Create(&newTransaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to create new transaction due to server error",
			Data:    nil,
		}
	}

	for _, detail := range transaction.Food {
		foodDetail := models.Transaction_Food{
			TransactionId: newTransaction.ID,
			FoodId:        detail.FoodId,
			Quantity:      detail.Quantity,
		}
		err := t.DB.Create(&foodDetail).Error
		if err != nil {
			return response.WebResponse{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to create new transaction detail due to server error",
				Data:    nil,
			}
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully create new transaction",
		Data:    nil,
	}
}

// DeleteTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) DeleteTransaction(TransactionId uint) response.WebResponse {
	var deletedTransaction models.Transaction
	err := t.DB.Where("id = ?", TransactionId).Find(&deletedTransaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction due to server error",
			Data:    nil,
		}
	} else if deletedTransaction.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Transaction not found",
			Data:    nil,
		}
	} else if deletedTransaction.AlreadyPay {
		return response.WebResponse{
			Code:    http.StatusNotModified,
			Status:  "Not modified",
			Message: "Can't delete paid transaction",
			Data:    nil,
		}
	}

	var detailTransaction []models.Transaction_Food
	res := t.DB.Where("transaction_id = ?", TransactionId).Delete(&detailTransaction)
	if res.Error != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction detail due to server error",
			Data:    nil,
		}
	} else if res.RowsAffected == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Transaction detail not found",
			Data:    nil,
		}
	}

	err = t.DB.Delete(&deletedTransaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to delete transaction due to server error",
			Data:    nil,
		}
	}

	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully delete transaction",
		Data:    nil,
	}
}

// GetAllPayTransaction implements TransacionRepositories.
func (t *TransactionRepositoriesImpl) GetAllPaidTransaction() response.WebResponse {
	transactions := []models.Transaction{}
	err := t.DB.Preload("Customer").Preload("Food").Where("already_pay = ?", true).Find(&transactions).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get all paid transaction due to server error",
			Data:    nil,
		}
	} else if len(transactions) == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Paid transaction not found",
			Data:    transactions,
		}
	}

	paidTransaction := []response.TransactionResponse{}
	for _, value := range transactions {
		transaction := response.TransactionResponse{
			TransactionId: value.ID,
			Customer: response.CustomerResponse{
				Id:           value.CustomerId,
				CustomerName: value.Customer.CustomerName,
			},
			AlreadyPay: value.AlreadyPay,
		}
		for _, food := range value.Food {
			var data models.Transaction_Food
			t.DB.Where("food_id = ? and transaction_id = ?", food.ID, value.ID).Find(&data)
			if data.FoodId == 0 && data.TransactionId == 0 {
				return response.WebResponse{
					Code:    http.StatusNotFound,
					Status:  "Not Found",
					Message: "Transaction with id " + strconv.FormatUint(uint64(data.TransactionId), 10) + " and food with id " + strconv.FormatUint(uint64(data.FoodId), 10) + " not found",
					Data:    nil,
				}
			}
			transaction.Food = append(transaction.Food, response.FoodDetailResponse{
				FoodId:    food.ID,
				FoodName:  food.FoodName,
				FoodPrice: food.FoodPrice,
				Quantity:  data.Quantity,
			})
		}
		paidTransaction = append(paidTransaction, transaction)
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get all paid transaction",
		Data:    paidTransaction,
	}
}

// GetTransactionById implements TransacionRepositories.
func (t *TransactionRepositoriesImpl) GetTransactionById(TransactionId uint) (response.WebResponse, bool) {
	var transaction models.Transaction
	err := t.DB.Preload("Customer").Preload("Food").Where("id = ?", TransactionId).Find(&transaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction due to server error",
			Data:    nil,
		}, false
	} else if transaction.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Transaction not found",
			Data:    nil,
		}, false
	}
	res := response.TransactionResponse{
		TransactionId: transaction.ID,
		Customer: response.CustomerResponse{
			Id:           transaction.Customer.ID,
			CustomerName: transaction.Customer.CustomerName,
		},
		AlreadyPay: transaction.AlreadyPay,
	}
	for _, foodDetail := range transaction.Food {
		var data models.Transaction_Food
		t.DB.Where("food_id = ? and transaction_id = ?", foodDetail.ID, transaction.ID).Find(&data)
		if data.FoodId == 0 && data.Transaction.ID == 0 {
			return response.WebResponse{
				Code:    http.StatusNotFound,
				Status:  "Not Found",
				Message: "Transaction with id " + strconv.FormatUint(uint64(data.TransactionId), 10) + " and food with id " + strconv.FormatUint(uint64(data.FoodId), 10) + " not found",
				Data:    nil,
			}, false
		}
		res.Food = append(res.Food, response.FoodDetailResponse{
			FoodId:    foodDetail.ID,
			FoodName:  foodDetail.FoodName,
			FoodPrice: foodDetail.FoodPrice,
			Quantity:  data.Quantity,
		})
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get transaction data",
		Data:    res,
	}, true
}

// GetUnpaidTransaction implements TransacionRepositories.
func (t *TransactionRepositoriesImpl) GetUnpaidTransaction() response.WebResponse {
	transactions := []models.Transaction{}
	err := t.DB.Preload("Customer").Preload("Food").Where("already_pay = ?", false).Find(&transactions).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get all unpaid transaction due to server error",
			Data:    nil,
		}
	} else if len(transactions) == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Unpaid transaction not found",
			Data:    transactions,
		}
	}

	unpaidTransaction := []response.TransactionResponse{}
	for _, value := range transactions {
		transaction := response.TransactionResponse{
			TransactionId: value.ID,
			Customer: response.CustomerResponse{
				Id:           value.CustomerId,
				CustomerName: value.Customer.CustomerName,
			},
			AlreadyPay: value.AlreadyPay,
		}
		for _, food := range value.Food {
			var data models.Transaction_Food
			t.DB.Where("food_id = ? and transaction_id = ?", food.ID, value.ID).Find(&data)
			if data.FoodId == 0 && data.TransactionId == 0 {
				return response.WebResponse{
					Code:    http.StatusNotFound,
					Status:  "Not Found",
					Message: "Transaction with id " + strconv.FormatUint(uint64(data.TransactionId), 10) + " and food with id " + strconv.FormatUint(uint64(data.FoodId), 10) + " not found",
					Data:    nil,
				}
			}
			transaction.Food = append(transaction.Food, response.FoodDetailResponse{
				FoodId:    food.ID,
				FoodName:  food.FoodName,
				FoodPrice: food.FoodPrice,
				Quantity:  data.Quantity,
			})
		}
		unpaidTransaction = append(unpaidTransaction, transaction)
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get all unpaid data",
		Data:    unpaidTransaction,
	}
}

// UpdateTransaction implements TransacionRepositories.
func (t *TransactionRepositoriesImpl) UpdateTransaction(transaction request.UpdateTransactionReq) response.WebResponse {
	var updatedTransaction models.Transaction
	err := t.DB.Where("id = ?", transaction.Id).Find(&updatedTransaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction due to server error",
			Data:    nil,
		}
	} else if updatedTransaction.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Transaction not found",
			Data:    nil,
		}
	} else if updatedTransaction.AlreadyPay {
		return response.WebResponse{
			Code:    http.StatusNotModified,
			Status:  "Not Modified",
			Message: "Can't delete paid transaction",
			Data:    nil,
		}
	}

	updatedTransaction.CustomerId = transaction.CustomerId
	err = t.DB.Updates(&updatedTransaction).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to update transaction due to server error",
			Data:    nil,
		}
	}

	var prevData []models.Transaction_Food
	res := t.DB.Where("transaction_id = ?", transaction.Id).Delete(&prevData)
	if res.Error != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get transaction detail due to server error",
			Data:    nil,
		}
	} else if res.RowsAffected == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Failed to delete old transaction detail due to server error",
			Data:    nil,
		}
	}

	for _, detail := range transaction.Food {
		var food models.Food
		err := t.DB.Where("id = ?", detail.FoodId).Find(&food).Error
		if err != nil {
			return response.WebResponse{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Can't get food due to server error",
				Data:    nil,
			}
		} else if food.ID == 0 {
			return response.WebResponse{
				Code:    http.StatusNotFound,
				Status:  "Not Found",
				Message: "Food with id " + strconv.FormatUint(uint64(detail.FoodId), 10) + " not found",
				Data:    nil,
			}
		}

		transactionFood := models.Transaction_Food{
			TransactionId: transaction.Id,
			FoodId:        detail.FoodId,
			Quantity:      detail.Quantity,
		}

		err = t.DB.Create(&transactionFood).Error
		if err != nil {
			return response.WebResponse{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to create new transaction detail due to server error",
				Data:    nil,
			}
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully update transaction",
		Data:    nil,
	}
}
