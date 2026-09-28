package models

type Transaction_Food struct {
	TransactionId uint `gorm:"primary_key" json:"TransactionId"`
	FoodId        uint `gorm:"primary_key" json:"FoodId"`
	Quantity      int  `json:"Quantity"`
	// UnitPrice adalah harga makanan pada saat transaksi dibuat.
	//
	// Tanpa kolom ini nilai transaksi lama ikut berubah setiap harga menu
	// diubah — struk kemarin berubah sendiri hari ini. Harga yang mengikat
	// adalah harga saat pesanan dibuat, jadi ia disalin, bukan dibaca ulang
	// dari tabel foods.
	UnitPrice   float64     `gorm:"not null;default:0" json:"UnitPrice"`
	Transaction Transaction `gorm:"foreignKey:TransactionId" json:"Transaction"`
	Food        Food        `gorm:"foreignKey:FoodId" json:"Food"`
}
