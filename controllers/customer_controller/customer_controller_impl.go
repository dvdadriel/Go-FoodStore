package customercontroller

import (
	"go-food-store/exception"
	"go-food-store/helpers"
	"go-food-store/json/request"
	customer "go-food-store/services/customer_service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CustomerControllerImpl struct {
	CustService customer.CustService
}

func NewCustomerController(CustService customer.CustService) CustomerController {
	return &CustomerControllerImpl{
		CustService: CustService,
	}
}

// CreateFood implements FoodController.
func (c *CustomerControllerImpl) CreateCust(w http.ResponseWriter, r *http.Request) {
	CustReq := &request.CreateCustReq{}
	err := helpers.Ummarshal(r, CustReq)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := c.CustService.Create(*CustReq)
		helpers.WriteJSON(w, response)
	}
}

// DeleteFood implements FoodController.
func (c *CustomerControllerImpl) DeleteCust(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id, err := strconv.ParseInt(param["custId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := c.CustService.Delete(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// GetAllFood implements FoodController.
func (c *CustomerControllerImpl) GetAllCust(w http.ResponseWriter, r *http.Request) {
	response := c.CustService.FindAll()
	helpers.WriteJSON(w, response)
}

// GetFoodById implements FoodController.
func (c *CustomerControllerImpl) GetCustById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["custId"], 0, 0)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		response := c.CustService.FindById(uint(id))
		helpers.WriteJSON(w, response)
	}
}

// UpdateFood implements FoodController.
func (c *CustomerControllerImpl) UpdateCust(w http.ResponseWriter, r *http.Request) {
	CustReq := &request.UpdateCustReq{}
	err := helpers.Ummarshal(r, CustReq)
	if err != nil {
		exception.BadRequestErr(w, r)
	} else {
		params := mux.Vars(r)
		id, err := strconv.ParseInt(params["custId"], 0, 0)
		if err != nil {
			exception.BadRequestErr(w, r)
		} else {
			CustReq.Id = uint(id)
			response := c.CustService.Update(*CustReq)
			helpers.WriteJSON(w, response)
		}
	}
}
