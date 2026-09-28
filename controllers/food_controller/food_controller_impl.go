package foodcontroller

import (
	"go-food-store/exception"
	"go-food-store/helpers"
	"go-food-store/json/request"
	food "go-food-store/services/food_service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FoodControllerImpl struct {
	FoodService food.FoodService
}

func NewFoodController(FoodService food.FoodService) FoodController {
	return &FoodControllerImpl{
		FoodService: FoodService,
	}
}

// CreateFood implements FoodController.
func (f *FoodControllerImpl) CreateFood(w http.ResponseWriter, r *http.Request) {
	FoodReq := &request.CreateFoodReq{}
	err := helpers.Ummarshal(r, FoodReq)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := f.FoodService.Create(*FoodReq)
		helpers.WriteJSON(w, response)
	}
}

// DeleteFood implements FoodController.
func (f *FoodControllerImpl) DeleteFood(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id, err := strconv.ParseInt(param["foodId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := f.FoodService.Delete(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// GetAllFood implements FoodController.
func (f *FoodControllerImpl) GetAllFood(w http.ResponseWriter, r *http.Request) {
	response := f.FoodService.FindAll(request.PageFrom(r))
	helpers.WriteJSON(w, response)
}

// GetFoodById implements FoodController.
func (f *FoodControllerImpl) GetFoodById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["foodId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := f.FoodService.FindById(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// UpdateFood implements FoodController.
func (f *FoodControllerImpl) UpdateFood(w http.ResponseWriter, r *http.Request) {
	FoodReq := &request.UpdateFoodReq{}
	err := helpers.Ummarshal(r, FoodReq)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		params := mux.Vars(r)
		id, err := strconv.ParseInt(params["foodId"], 0, 0)
		if err != nil {
			exception.BadRequestErr(w, r)
		} else {
			FoodReq.Id = uint(id)
			response := f.FoodService.Update(*FoodReq)
			helpers.WriteJSON(w, response)
		}
	}
}
