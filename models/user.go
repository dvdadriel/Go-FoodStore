package models

import "gorm.io/gorm"

// Peran yang dikenal. Sengaja hanya dua: satu yang mengelola toko, satu yang
// memesan. Peran ketiga ditambahkan kalau memang ada pekerjaan yang tidak
// muat di keduanya, bukan sebelumnya.
const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100;uniqueIndex;not null" json:"Username"`
	// PasswordHash, bukan Password: yang tersimpan adalah hasil PBKDF2
	// berisi salt dan jumlah iterasinya. Tanpa tag json supaya tidak
	// mungkin ikut terserialisasi ke response.
	PasswordHash string `gorm:"not null" json:"-"`
	Role         string `gorm:"size:20;not null;default:customer" json:"Role"`
}
