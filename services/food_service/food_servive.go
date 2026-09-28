package foodservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
)

type FoodService interface {
	Create(food request.CreateFoodReq) response.WebResponse
	Update(food request.UpdateFoodReq) response.WebResponse
	Delete(FoodId uint) response.WebResponse
	FindAll(page request.Page) response.WebResponse
	FindById(FoodId uint) response.WebResponse
}
