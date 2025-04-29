package customerrepository

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	"net/http"

	"gorm.io/gorm"
)

type CustomerRepoImpl struct {
	DB *gorm.DB
}

func NewCustomerRepo(DB *gorm.DB) CustomerRepo {
	return &CustomerRepoImpl{
		DB: DB,
	}
}

// CreateCustomer implements CustomerRepo.
func (c *CustomerRepoImpl) CreateCustomer(cust models.Customer) response.WebResponse {
	err := c.DB.Create(&cust).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't create customer due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully create new customer",
		Data:    nil,
	}
}

// DeleteCustomer implements CustomerRepo.
func (c *CustomerRepoImpl) DeleteCustomer(CustId uint) response.WebResponse {
	var customer models.Customer
	err := c.DB.Where("id = ?", CustId).Delete(&customer).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't delete customer due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully delete customer",
		Data:    nil,
	}
}

// GetAllCustomer implements CustomerRepo.
func (c *CustomerRepoImpl) GetAllCustomer() response.WebResponse {
	var customers []models.Customer
	err := c.DB.Find(&customers).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get all customer due to server error",
			Data:    nil,
		}
	}
	customerResponse := []response.CustomerResponse{}
	for _, value := range customers {
		customerResponse = append(customerResponse, response.CustomerResponse{
			Id:           value.ID,
			CustomerName: value.CustomerName,
		})
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get all customer",
		Data:    customerResponse,
	}
}

// GetCustomerById implements CustomerRepo.
func (c *CustomerRepoImpl) GetCustomerById(CustId uint) (response.WebResponse, bool) {
	var customer models.Customer
	err := c.DB.Where("id = ?", CustId).Find(&customer).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't get customer due to server error",
			Data:    nil,
		}, false
	} else if customer.ID == 0 {
		return response.WebResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: "Customer not found",
			Data:    nil,
		}, false
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully get customer data",
		Data: response.CustomerResponse{
			Id:           customer.ID,
			CustomerName: customer.CustomerName,
		},
	}, true
}

// UpdateCustomer implements CustomerRepo.
func (c *CustomerRepoImpl) UpdateCustomer(cust models.Customer) response.WebResponse {
	updatedCust := request.UpdateCustReq{
		Id:           cust.ID,
		CustomerName: cust.CustomerName,
	}
	err := c.DB.Model(&cust).Updates(updatedCust).Error
	if err != nil {
		return response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Can't update customer due to server error",
			Data:    nil,
		}
	}
	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully update customer data",
		Data:    nil,
	}
}
