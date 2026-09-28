package response

type TransactionResponse struct {
	TransactionId uint
	Customer      CustomerResponse
	Food          []FoodDetailResponse
	TotalPrice    float64
	AlreadyPay    bool
}
