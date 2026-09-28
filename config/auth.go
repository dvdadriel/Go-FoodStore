package config

import (
	"log/slog"
	"strconv"
	"time"

	"gorm.io/gorm"

	"go-food-store/auth"
	"go-food-store/models"
	userrepository "go-food-store/repositories/user_repository"
)

// JWTSecret adalah kunci penanda tangan token.
//
// Tidak ada nilai default yang aman untuk ini, jadi ketiadaannya dijadikan
// kegagalan yang berisik saat start — bukan diganti kunci acak (token jadi
// tidak berlaku setiap restart dan antar instance) dan bukan diganti kunci
// tetap yang bisa dibaca siapa pun di repo ini.
func JWTSecret() string {
	secret := env("JWT_SECRET", "")
	if len(secret) < 32 {
		panic("JWT_SECRET wajib di-set dan minimal 32 karakter — lihat .env.example")
	}
	return secret
}

// TokenTTL adalah masa berlaku token akses, dalam detik.
func TokenTTL() time.Duration {
	raw := env("TOKEN_TTL_SECONDS", "3600")
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		slog.Warn("TOKEN_TTL_SECONDS tidak sah, memakai 3600 detik", "nilai", raw)
		seconds = 3600
	}
	return time.Duration(seconds) * time.Second
}

// SeedAdmin membuat akun admin pertama dari environment kalau belum ada.
//
// Admin tidak bisa lahir dari endpoint pendaftaran — kalau bisa, siapa pun
// tinggal mendaftar sebagai admin. Jadi ia harus datang dari luar aplikasi,
// dan environment adalah jalan yang paling sederhana.
//
// Kalau ADMIN_USERNAME atau ADMIN_PASSWORD tidak di-set, langkah ini
// dilewati: instalasi tanpa admin lebih baik daripada instalasi dengan admin
// berkata sandi yang bisa ditebak dari kode sumber.
func SeedAdmin(db *gorm.DB) {
	username := env("ADMIN_USERNAME", "")
	password := env("ADMIN_PASSWORD", "")
	if username == "" || password == "" {
		slog.Warn("ADMIN_USERNAME/ADMIN_PASSWORD tidak di-set, akun admin tidak dibuat")
		return
	}

	repo := userrepository.NewUserRepo(db)
	exists, err := repo.ExistsByUsername(username)
	if err != nil {
		slog.Error("gagal memeriksa akun admin", "error", err)
		return
	}
	if exists {
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		slog.Error("gagal membuat hash kata sandi admin", "error", err)
		return
	}
	if _, err := repo.Create(models.User{
		Username:     username,
		PasswordHash: hash,
		Role:         models.RoleAdmin,
	}); err != nil {
		slog.Error("gagal membuat akun admin", "error", err)
		return
	}
	slog.Info("akun admin dibuat", "username", username)
}
