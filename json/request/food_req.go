package request

type CreateFoodReq struct {
	FoodName string `validate:"required,min=1,max=200" json:"FoodName"`
	// gt=0, bukan required: `required` menolak 0 tapi meloloskan harga
	// negatif, karena yang diperiksanya hanya "bukan nilai nol".
	FoodPrice float64 `validate:"gt=0" json:"FoodPrice"`
}

type UpdateFoodReq struct {
	Id        uint    `validate:"required" json:"Id"`
	FoodName  string  `validate:"required,min=1,max=200" json:"FoodName"`
	FoodPrice float64 `validate:"gt=0" json:"FoodPrice"`
}

type FoodDetailReq struct {
	FoodId uint `validate:"required" json:"FoodId"`
	// Sebelumnya tanpa aturan sama sekali: quantity 0 membuat baris pesanan
	// yang tidak berarti, dan quantity negatif mengurangi total belanja.
	Quantity int `validate:"gt=0" json:"Quantity"`
}
