package exception

import (
	"go-food-store/helpers"
	"go-food-store/json/response"
	"net/http"
)

func NotFoundErr(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, response.WebResponse{
		Code:    http.StatusNotFound,
		Status:  "Not Found",
		Message: "Data not found",
		Data:    nil,
	})
}

func NotAllowedErr(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, response.WebResponse{
		Code:    http.StatusMethodNotAllowed,
		Status:  "Method Not Allowed",
		Message: "Wrong method",
		Data:    nil,
	})
}

func BadRequestErr(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, response.WebResponse{
		Code:    http.StatusBadRequest,
		Status:  "Bad Request",
		Message: "Please check the request",
		Data:    nil,
	})
}
