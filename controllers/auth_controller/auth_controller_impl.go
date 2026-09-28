package authcontroller

import (
	"net/http"

	"go-food-store/exception"
	"go-food-store/helpers"
	"go-food-store/json/request"
	authservice "go-food-store/services/auth_service"
)

type AuthControllerImpl struct {
	AuthService authservice.AuthService
}

func NewAuthController(service authservice.AuthService) AuthController {
	return &AuthControllerImpl{AuthService: service}
}

func (a *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	req := &request.RegisterReq{}
	if err := helpers.Ummarshal(r, req); err != nil {
		exception.BadRequestErr(w, r)
		return
	}
	helpers.WriteJSON(w, a.AuthService.Register(*req))
}

func (a *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	req := &request.LoginReq{}
	if err := helpers.Ummarshal(r, req); err != nil {
		exception.BadRequestErr(w, r)
		return
	}
	helpers.WriteJSON(w, a.AuthService.Login(*req))
}
