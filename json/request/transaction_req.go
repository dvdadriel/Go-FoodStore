package request

type CreateTransactionReq struct {
	CustomerId uint            `validate:"required"`
	Food       []FoodDetailReq `validate:"required"`
}

type UpdateTransactionReq struct {
	Id         uint            `validate:"required"`
	CustomerId uint            `validate:"required"`
	Food       []FoodDetailReq `validate:"required"`
}
