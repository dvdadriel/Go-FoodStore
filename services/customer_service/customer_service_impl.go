package customerservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	customer "go-food-store/repositories/customer_repository"
	"net/http"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type CustServiceImpl struct {
	CustRepo  customer.CustomerRepo
	Validator *validator.Validate
}

func NewCustService(Repo customer.CustomerRepo, Validator *validator.Validate) CustService {
	return &CustServiceImpl{
		CustRepo:  Repo,
		Validator: Validator,
	}
}

// Create implements CustService.
func (c *CustServiceImpl) Create(cust request.CreateCustReq) response.WebResponse {
	err := c.Validator.Struct(&cust)
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Please check your request",
			Data:    nil,
		}
	}

	newCust := models.Customer{
		CustomerName: cust.CustomerName,
	}
	response := c.CustRepo.CreateCustomer(newCust)
	return response
}

// Delete implements CustService.
func (c *CustServiceImpl) Delete(CustId uint) response.WebResponse {
	response, found := c.CustRepo.GetCustomerById(CustId)
	if !found {
		return response
	}
	response = c.CustRepo.DeleteCustomer(CustId)
	return response
}

// FindAll implements CustService.
func (c *CustServiceImpl) FindAll() response.WebResponse {
	response := c.CustRepo.GetAllCustomer()
	return response
}

// FindById implements CustService.
func (c *CustServiceImpl) FindById(CustId uint) response.WebResponse {
	response, _ := c.CustRepo.GetCustomerById(CustId)
	return response
}

// Update implements CustService.
func (c *CustServiceImpl) Update(cust request.UpdateCustReq) response.WebResponse {
	response, found := c.CustRepo.GetCustomerById(cust.Id)
	if !found {
		return response
	}
	updatedCust := models.Customer{
		Model: gorm.Model{
			ID: cust.Id,
		},
		CustomerName: cust.CustomerName,
	}
	response = c.CustRepo.UpdateCustomer(updatedCust)
	return response
}
