package request

type CreateFoodReq struct {
	FoodName  string  `validate:"required,min=1,max=200" json:"FoodName"`
	FoodPrice float64 `validate:"required"`
}

type UpdateFoodReq struct {
	Id        uint    `validate:"required"`
	FoodName  string  `validate:"required,min=1,max=200" json:"FoodName"`
	FoodPrice float64 `validate:"required"`
}

type FoodDetailReq struct {
	FoodId   uint
	Quantity int
}
