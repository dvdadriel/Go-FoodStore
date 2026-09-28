package foodrepository

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
)

type FoodRepo interface {
	CreateFood(food models.Food) response.WebResponse
	UpdateFood(food models.Food) response.WebResponse
	DeleteFood(FoodId uint) response.WebResponse
	GetAllFood(page request.Page) response.WebResponse
	GetFoodById(FoodId uint) (response.WebResponse, bool)
}
