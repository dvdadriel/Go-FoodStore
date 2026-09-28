package foodrepository

import (
	"errors"

	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	"net/http"

	"gorm.io/gorm"
)

type FoodRepoImpl struct {
	DB *gorm.DB
}

func NewFoodRepo(DB *gorm.DB) FoodRepo {
	return &FoodRepoImpl{
		DB: DB,
	}
}

// CreateFood implements FoodRepo.
func (f *FoodRepoImpl) CreateFood(food models.Food) response.WebResponse {
	err := f.DB.Create(&food).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't create new food due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully create new food",
		Data:    nil,
	}
}

// DeleteFood implements FoodRepo.
func (f *FoodRepoImpl) DeleteFood(FoodId uint) response.WebResponse {
	var food models.Food
	err := f.DB.Where("id = ?", FoodId).Delete(&food).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't delete food due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully delete food",
		Data:    nil,
	}
}

// GetAllFood implements FoodRepo.
func (f *FoodRepoImpl) GetAllFood() response.WebResponse {
	foods := []models.Food{}
	err := f.DB.Find(&foods).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get food due to server error",
			Data:    foods,
		}
	}
	foodResponse := []response.FoodResponse{}
	for _, value := range foods {
		foodResponse = append(foodResponse, response.FoodResponse{
			Id:        value.ID,
			FoodName:  value.FoodName,
			FoodPrice: value.FoodPrice,
		})
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get all food",
		Data:    foodResponse,
	}
}

// GetFoodById implements FoodRepo.
func (f *FoodRepoImpl) GetFoodById(FoodId uint) (response.WebResponse, bool) {
	var food models.Food
	// First, bukan Find: Find pada satu struct membalas nil error untuk nol
	// baris, sehingga kegagalan database yang sebenarnya tidak terbedakan
	// dari data yang memang tidak ada.
	err := f.DB.First(&food, FoodId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Food not found",
			Data:    nil,
		}, false
	} else if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get food due to server error",
			Data:    nil,
		}, false
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get food data",
		Data: response.FoodResponse{
			Id:        food.ID,
			FoodName:  food.FoodName,
			FoodPrice: food.FoodPrice,
		},
	}, true
}

// UpdateFood implements FoodRepo.
func (f *FoodRepoImpl) UpdateFood(food models.Food) response.WebResponse {
	updatedFood := request.UpdateFoodReq{
		Id:        food.ID,
		FoodName:  food.FoodName,
		FoodPrice: food.FoodPrice,
	}
	err := f.DB.Model(&food).Updates(updatedFood).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't update food due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully update food data",
		Data:    nil,
	}
}
