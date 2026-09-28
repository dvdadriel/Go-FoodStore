package foodservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	food "go-food-store/repositories/food_repository"
	"net/http"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type FoodServiceImpl struct {
	FoodRepo  food.FoodRepo
	Validator *validator.Validate
}

func NewFoodService(Repo food.FoodRepo, Validator *validator.Validate) FoodService {
	return &FoodServiceImpl{
		FoodRepo:  Repo,
		Validator: Validator,
	}
}

// Create implements FoodService.
func (f *FoodServiceImpl) Create(food request.CreateFoodReq) response.WebResponse {
	err := f.Validator.Struct(food)
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Please check the request",
			Data:    nil,
		}
	}

	Newfood := models.Food{
		FoodName:  food.FoodName,
		FoodPrice: food.FoodPrice,
	}
	response := f.FoodRepo.CreateFood(Newfood)
	return response
}

// Delete implements FoodService.
func (f *FoodServiceImpl) Delete(FoodId uint) response.WebResponse {
	response, found := f.FoodRepo.GetFoodById(FoodId)
	if !found {
		return response
	}
	response = f.FoodRepo.DeleteFood(FoodId)
	return response
}

// FindAll implements FoodService.
func (f *FoodServiceImpl) FindAll(page request.Page) response.WebResponse {
	return f.FoodRepo.GetAllFood(page)
}

// FindById implements FoodService.
func (f *FoodServiceImpl) FindById(FoodId uint) response.WebResponse {
	response, _ := f.FoodRepo.GetFoodById(FoodId)
	return response
}

// Update implements FoodService.
func (f *FoodServiceImpl) Update(food request.UpdateFoodReq) response.WebResponse {
	err := f.Validator.Struct(food)
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Please check your request",
			Data:    nil,
		}
	}
	response, found := f.FoodRepo.GetFoodById(food.Id)
	if !found {
		return response
	}
	updatefood := models.Food{
		Model: gorm.Model{
			ID: food.Id,
		},
		FoodName:  food.FoodName,
		FoodPrice: food.FoodPrice,
	}
	response = f.FoodRepo.UpdateFood(updatefood)
	return response
}
