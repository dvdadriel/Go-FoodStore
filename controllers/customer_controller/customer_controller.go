package customercontroller

import "net/http"

type CustomerController interface {
	CreateCust(w http.ResponseWriter, r *http.Request)
	UpdateCust(w http.ResponseWriter, r *http.Request)
	DeleteCust(w http.ResponseWriter, r *http.Request)
	GetCustById(w http.ResponseWriter, r *http.Request)
	GetAllCust(w http.ResponseWriter, r *http.Request)
}
