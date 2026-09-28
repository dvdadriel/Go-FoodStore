package response

type FoodResponse struct {
	Id        uint
	FoodName  string
	FoodPrice float64
}

type FoodDetailResponse struct {
	FoodId   uint
	FoodName string
	// FoodPrice adalah harga yang tercatat saat transaksi dibuat, bukan
	// harga menu saat ini.
	FoodPrice float64
	Quantity  int
	Subtotal  float64
}
