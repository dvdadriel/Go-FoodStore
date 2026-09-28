package authservice

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"gorm.io/gorm"

	"go-food-store/auth"
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
	userrepository "go-food-store/repositories/user_repository"
)

type AuthServiceImpl struct {
	UserRepo  userrepository.UserRepo
	Issuer    *auth.Issuer
	TokenTTL  int64
	Validator *validator.Validate
}

func NewAuthService(repo userrepository.UserRepo, issuer *auth.Issuer, ttlSeconds int64, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepo:  repo,
		Issuer:    issuer,
		TokenTTL:  ttlSeconds,
		Validator: validate,
	}
}

func webError(code int, status, message string) response.WebResponse {
	return response.WebResponse{Code: code, Status: status, Message: message, Data: nil}
}

// Register selalu membuat akun dengan peran customer.
//
// Peran admin tidak pernah bisa diminta lewat pendaftaran terbuka — kalau
// bisa, siapa pun tinggal mendaftar sebagai admin. Admin dibuat saat startup
// dari environment (lihat config.SeedAdmin).
func (a *AuthServiceImpl) Register(req request.RegisterReq) response.WebResponse {
	if err := a.Validator.Struct(req); err != nil {
		return webError(http.StatusBadRequest, "Bad Request", "Username minimal 3 karakter dan password minimal 8 karakter")
	}

	exists, err := a.UserRepo.ExistsByUsername(req.Username)
	if err != nil {
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't register due to server error")
	}
	if exists {
		return webError(http.StatusConflict, "Conflict", "Username already taken")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't register due to server error")
	}

	created, err := a.UserRepo.Create(models.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         models.RoleCustomer,
	})
	if err != nil {
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't register due to server error")
	}

	return response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Successfully register new user",
		Data: response.UserResponse{
			Id:       created.ID,
			Username: created.Username,
			Role:     created.Role,
		},
	}
}

// Login menukar kredensial dengan token akses.
func (a *AuthServiceImpl) Login(req request.LoginReq) response.WebResponse {
	if err := a.Validator.Struct(req); err != nil {
		return webError(http.StatusBadRequest, "Bad Request", "Please check the request")
	}

	// Pesan yang sama untuk username tidak ada dan password salah. Pesan
	// yang berbeda akan memberi tahu penyerang username mana yang terdaftar.
	const gagal = "Invalid username or password"

	user, err := a.UserRepo.FindByUsername(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return webError(http.StatusUnauthorized, "Unauthorized", gagal)
	} else if err != nil {
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't login due to server error")
	}

	cocok, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		// Hash yang tidak terbaca adalah data rusak, bukan kesalahan
		// pengguna — dicatat supaya tidak lenyap menjadi "password salah".
		slog.Error("hash kata sandi tidak terbaca", "user_id", user.ID, "error", err)
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't login due to server error")
	}
	if !cocok {
		return webError(http.StatusUnauthorized, "Unauthorized", gagal)
	}

	token, err := a.Issuer.Issue(user.ID, user.Username, user.Role)
	if err != nil {
		return webError(http.StatusInternalServerError, "Internal Server Error", "Can't issue token due to server error")
	}

	return response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Successfully login",
		Data: response.LoginResponse{
			Token:     token,
			ExpiresIn: a.TokenTTL,
			User: response.UserResponse{
				Id:       user.ID,
				Username: user.Username,
				Role:     user.Role,
			},
		},
	}
}
