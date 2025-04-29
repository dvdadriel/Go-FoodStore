package exception

import (
	"encoding/json"
	"go-food-store/helpers"
	"go-food-store/json/response"
	"net/http"
)

func NotFoundErr(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(response.WebResponse{
		Code:    http.StatusNotFound,
		Status:  "Not Found",
		Message: "Data not found",
		Data:    nil,
	})
	helpers.PanicHelper(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func NotAllowedErr(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(response.WebResponse{
		Code:    http.StatusMethodNotAllowed,
		Status:  "Method Not Allowed",
		Message: "Wrong method",
		Data:    nil,
	})
	helpers.PanicHelper(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func BadRequestErr(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(response.WebResponse{
		Code:    http.StatusBadRequest,
		Status:  "Bad Request",
		Message: "Please check the request",
		Data:    nil,
	})
	helpers.PanicHelper(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
