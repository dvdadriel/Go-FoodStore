package customerservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
)

type CustService interface {
	Create(cust request.CreateCustReq) response.WebResponse
	Update(cust request.UpdateCustReq) response.WebResponse
	Delete(CustId uint) response.WebResponse
	FindAll() response.WebResponse
	FindById(CustId uint) response.WebResponse
}
