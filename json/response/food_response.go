package response

type FoodResponse struct {
	Id        uint
	FoodName  string
	FoodPrice float64
}

type FoodDetailResponse struct {
	FoodId    uint
	FoodName  string
	FoodPrice float64
	Quantity  int
}
