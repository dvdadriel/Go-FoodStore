package response

type TransactionResponse struct {
	TransactionId uint
	Customer      CustomerResponse
	Food          []FoodDetailResponse
	AlreadyPay    bool
}
