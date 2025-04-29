package customerrepository

import (
	"go-food-store/json/response"
	"go-food-store/models"
)

type CustomerRepo interface {
	CreateCustomer(cust models.Customer) response.WebResponse
	UpdateCustomer(cust models.Customer) response.WebResponse
	DeleteCustomer(CustId uint) response.WebResponse
	GetAllCustomer() response.WebResponse
	GetCustomerById(CustId uint) (response.WebResponse, bool)
}
