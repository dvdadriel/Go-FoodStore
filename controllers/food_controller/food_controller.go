package foodcontroller

import "net/http"

type FoodController interface {
	CreateFood(w http.ResponseWriter, r *http.Request)
	UpdateFood(w http.ResponseWriter, r *http.Request)
	DeleteFood(w http.ResponseWriter, r *http.Request)
	GetFoodById(w http.ResponseWriter, r *http.Request)
	GetAllFood(w http.ResponseWriter, r *http.Request)
}
