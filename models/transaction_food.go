package models

type Transaction_Food struct {
	TransactionId uint        `gorm:"primary_key" json:"TransactionId"`
	FoodId        uint        `gorm:"primary_key" json:"FoodId"`
	Quantity      int         `json:"Quantity"`
	Transaction   Transaction `gorm:"foreignKey:TransactionId" json:"Transaction"`
	Food          Food        `gorm:"foreignKey:FoodId" json:"Food"`
}
