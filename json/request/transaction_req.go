package request

type CreateTransactionReq struct {
	CustomerId uint `validate:"required" json:"CustomerId"`
	// dive meneruskan validasi ke tiap elemen. Tanpa itu isi slice tidak
	// pernah diperiksa dan quantity 0 lolos begitu saja.
	Food []FoodDetailReq `validate:"required,min=1,dive" json:"Food"`
}

type UpdateTransactionReq struct {
	Id         uint            `validate:"required" json:"Id"`
	CustomerId uint            `validate:"required" json:"CustomerId"`
	Food       []FoodDetailReq `validate:"required,min=1,dive" json:"Food"`
}
