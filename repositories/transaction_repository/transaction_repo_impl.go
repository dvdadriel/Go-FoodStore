package transaction

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
)

type TransactionRepositoriesImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepositories {
	return &TransactionRepositoriesImpl{
		DB: DB,
	}
}

func serverErr(message string) response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusInternalServerError,
		Status:  "Internal Server Error",
		Message: message,
		Data:    nil,
	}
}

func notFound(message string) response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusNotFound,
		Status:  "Not Found",
		Message: message,
		Data:    nil,
	}
}

func ok(message string, data interface{}) response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: message,
		Data:    data,
	}
}

func id(v uint) string { return strconv.FormatUint(uint64(v), 10) }

// GetCustomerById implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetCustomerById(CustomerId uint) (response.WebResponse, bool) {
	var customer models.Customer
	// First, bukan Find: Find pada satu struct membalas nil error untuk nol
	// baris, jadi "tidak ada datanya" harus ditebak dari ID == 0. First
	// membedakan baris kosong dari kegagalan database yang sebenarnya.
	err := t.DB.First(&customer, CustomerId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Customer with id " + id(CustomerId) + " not found"), false
	} else if err != nil {
		return serverErr("Can't get customer due to server error"), false
	}
	return response.WebResponse{}, true
}

// GetFoodById implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetFoodById(FoodId uint) (response.WebResponse, bool) {
	var food models.Food
	err := t.DB.First(&food, FoodId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Food with id " + id(FoodId) + " not found"), false
	} else if err != nil {
		return serverErr("Can't get food due to server error"), false
	}
	return response.WebResponse{}, true
}

// AcceptPayment implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) AcceptPayment(TransactionId uint) response.WebResponse {
	var transaction models.Transaction
	err := t.DB.First(&transaction, TransactionId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Transaction not found")
	} else if err != nil {
		return serverErr("Can't get transaction due to server error")
	}
	if transaction.AlreadyPay {
		return conflict("Transaction has already been paid")
	}
	if err := t.DB.Model(&transaction).Update("already_pay", true).Error; err != nil {
		return serverErr("Failed to update transaction status due to server error")
	}
	return ok("Successfully update transaction status", nil)
}

// conflict dipakai untuk permintaan yang sah bentuknya tapi bentrok dengan
// keadaan data — misalnya mengubah transaksi yang sudah lunas.
//
// Sebelumnya kasus ini membalas 304 Not Modified. 304 dilarang membawa body
// oleh spesifikasi HTTP, jadi pesan alasannya hilang di jalan dan client
// hanya melihat response kosong.
func conflict(message string) response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusConflict,
		Status:  "Conflict",
		Message: message,
		Data:    nil,
	}
}

// CreateTransaction implements TransactionRepositories.
//
// Seluruhnya di dalam satu transaksi database. Sebelumnya header di-insert
// lebih dulu dan tiap baris detail menyusul satu per satu; kegagalan di
// tengah meninggalkan transaksi yatim tanpa item, yang tetap terhitung
// sebagai pesanan belum bayar selamanya.
func (t *TransactionRepositoriesImpl) CreateTransaction(req request.CreateTransactionReq) response.WebResponse {
	var created models.Transaction
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		created = models.Transaction{
			CustomerId: req.CustomerId,
			AlreadyPay: false,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return insertDetails(tx, created.ID, req.Food)
	})
	if err != nil {
		if res, isNotFound := asNotFound(err); isNotFound {
			return res
		}
		return serverErr("Failed to create new transaction due to server error")
	}
	return ok("Successfully create new transaction", map[string]uint{"TransactionId": created.ID})
}

// foodNotFoundError membawa id makanan yang tidak ada keluar dari closure
// db.Transaction, supaya rollback tetap terjadi tapi client tetap dapat 404
// yang spesifik alih-alih 500 generik.
type foodNotFoundError struct{ foodId uint }

func (e foodNotFoundError) Error() string { return "food " + id(e.foodId) + " not found" }

func asNotFound(err error) (response.WebResponse, bool) {
	var e foodNotFoundError
	if errors.As(err, &e) {
		return notFound("Food with id " + id(e.foodId) + " not found"), true
	}
	return response.WebResponse{}, false
}

// insertDetails menyalin harga makanan yang berlaku sekarang ke tiap baris
// detail, lalu menyimpannya sekaligus.
func insertDetails(tx *gorm.DB, transactionId uint, items []request.FoodDetailReq) error {
	details := make([]models.Transaction_Food, 0, len(items))
	for _, item := range items {
		var food models.Food
		if err := tx.First(&food, item.FoodId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return foodNotFoundError{foodId: item.FoodId}
			}
			return err
		}
		details = append(details, models.Transaction_Food{
			TransactionId: transactionId,
			FoodId:        item.FoodId,
			Quantity:      item.Quantity,
			UnitPrice:     food.FoodPrice,
		})
	}
	if len(details) == 0 {
		return nil
	}
	return tx.Create(&details).Error
}

// DeleteTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) DeleteTransaction(TransactionId uint) response.WebResponse {
	var existing models.Transaction
	err := t.DB.First(&existing, TransactionId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Transaction not found")
	} else if err != nil {
		return serverErr("Can't get transaction due to server error")
	}
	if existing.AlreadyPay {
		return conflict("Can't delete paid transaction")
	}

	err = t.DB.Transaction(func(tx *gorm.DB) error {
		// Tanpa RowsAffected == 0 sebagai error: transaksi tanpa item itu
		// sah untuk dihapus, dan sebelumnya kasus itu membalas 404 padahal
		// transaksinya jelas ada.
		if err := tx.Where("transaction_id = ?", TransactionId).Delete(&models.Transaction_Food{}).Error; err != nil {
			return err
		}
		return tx.Delete(&existing).Error
	})
	if err != nil {
		return serverErr("Failed to delete transaction due to server error")
	}
	return ok("Successfully delete transaction", nil)
}

// UpdateTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) UpdateTransaction(req request.UpdateTransactionReq) response.WebResponse {
	var existing models.Transaction
	err := t.DB.First(&existing, req.Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Transaction not found")
	} else if err != nil {
		return serverErr("Can't get transaction due to server error")
	}
	if existing.AlreadyPay {
		return conflict("Can't update paid transaction")
	}

	err = t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Update("customer_id", req.CustomerId).Error; err != nil {
			return err
		}
		if err := tx.Where("transaction_id = ?", req.Id).Delete(&models.Transaction_Food{}).Error; err != nil {
			return err
		}
		return insertDetails(tx, req.Id, req.Food)
	})
	if err != nil {
		if res, isNotFound := asNotFound(err); isNotFound {
			return res
		}
		return serverErr("Failed to update transaction due to server error")
	}
	return ok("Successfully update transaction", nil)
}

// GetAllPaidTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetAllPaidTransaction(page request.Page) response.WebResponse {
	return t.listByPaymentStatus(page, true, "Successfully get all paid transaction")
}

// GetUnpaidTransaction implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetUnpaidTransaction(page request.Page) response.WebResponse {
	return t.listByPaymentStatus(page, false, "Successfully get all unpaid transaction")
}

// listByPaymentStatus mengambil daftar transaksi beserta itemnya.
//
// Daftar kosong membalas 200 dengan array kosong, bukan 404. Tidak adanya
// transaksi lunas bukan kesalahan — itu jawaban yang benar untuk toko yang
// baru buka, dan client tidak perlu memperlakukannya sebagai error.
func (t *TransactionRepositoriesImpl) listByPaymentStatus(page request.Page, paid bool, message string) response.WebResponse {
	var transactions []models.Transaction
	if err := t.DB.Preload("Customer").Where("already_pay = ?", paid).Limit(page.Limit).Offset(page.Offset).Order("id").Find(&transactions).Error; err != nil {
		return serverErr("Can't get transaction due to server error")
	}

	result, err := t.buildResponses(transactions)
	if err != nil {
		return serverErr("Can't get transaction detail due to server error")
	}
	return ok(message, result)
}

// GetTransactionById implements TransactionRepositories.
func (t *TransactionRepositoriesImpl) GetTransactionById(TransactionId uint) (response.WebResponse, bool) {
	var transaction models.Transaction
	err := t.DB.Preload("Customer").First(&transaction, TransactionId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound("Transaction not found"), false
	} else if err != nil {
		return serverErr("Can't get transaction due to server error"), false
	}

	result, err := t.buildResponses([]models.Transaction{transaction})
	if err != nil {
		return serverErr("Can't get transaction detail due to server error"), false
	}
	return ok("Successfully get transaction data", result[0]), true
}

// buildResponses mengambil seluruh detail dalam satu query, bukan satu query
// per item per transaksi seperti sebelumnya — 50 transaksi berisi 5 item
// dulunya berarti 250 perjalanan ke database.
func (t *TransactionRepositoriesImpl) buildResponses(transactions []models.Transaction) ([]response.TransactionResponse, error) {
	result := make([]response.TransactionResponse, 0, len(transactions))
	if len(transactions) == 0 {
		return result, nil
	}

	ids := make([]uint, 0, len(transactions))
	for _, tr := range transactions {
		ids = append(ids, tr.ID)
	}

	var details []models.Transaction_Food
	if err := t.DB.Preload("Food").Where("transaction_id IN ?", ids).Find(&details).Error; err != nil {
		return nil, err
	}

	byTransaction := map[uint][]models.Transaction_Food{}
	for _, d := range details {
		byTransaction[d.TransactionId] = append(byTransaction[d.TransactionId], d)
	}

	for _, tr := range transactions {
		item := response.TransactionResponse{
			TransactionId: tr.ID,
			Customer: response.CustomerResponse{
				Id:           tr.CustomerId,
				CustomerName: tr.Customer.CustomerName,
			},
			Food:       []response.FoodDetailResponse{},
			AlreadyPay: tr.AlreadyPay,
		}
		for _, d := range byTransaction[tr.ID] {
			subtotal := d.UnitPrice * float64(d.Quantity)
			item.Food = append(item.Food, response.FoodDetailResponse{
				FoodId:    d.FoodId,
				FoodName:  d.Food.FoodName,
				FoodPrice: d.UnitPrice,
				Quantity:  d.Quantity,
				Subtotal:  subtotal,
			})
			item.TotalPrice += subtotal
		}
		result = append(result, item)
	}
	return result, nil
}
