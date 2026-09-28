package authservice

import (
	"go-food-store/json/request"
	"go-food-store/json/response"
)

type AuthService interface {
	Register(req request.RegisterReq) response.WebResponse
	Login(req request.LoginReq) response.WebResponse
}
