// Package auth menangani kata sandi dan token akses.
//
// Seluruhnya memakai pustaka standar. crypto/pbkdf2 masuk stdlib di Go 1.24,
// dan HMAC-SHA256 sudah lama ada, jadi tidak ada alasan menambah dependensi
// untuk keduanya.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Parameter PBKDF2. iterasi dinaikkan dengan menaikkan konstanta ini; hash
// lama tetap bisa diverifikasi karena jumlah iterasinya ikut tersimpan di
// dalam string hash-nya.
const (
	pbkdf2Iterations = 210000
	saltLength       = 16
	keyLength        = 32
)

var ErrInvalidHash = errors.New("format hash kata sandi tidak dikenali")

// HashPassword mengembalikan hash yang aman disimpan, berisi salt dan jumlah
// iterasinya sekaligus: "pbkdf2-sha256$iterasi$salt$hash".
//
// Salt disimpan menyatu dengan hash-nya, bukan di kolom terpisah, supaya
// keduanya tidak mungkin terpisah saat baris disalin atau dimigrasikan.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iterations, keyLength)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("pbkdf2-sha256$%d$%s$%s",
		pbkdf2Iterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

// VerifyPassword memeriksa kata sandi terhadap hash tersimpan.
//
// Perbandingannya memakai subtle.ConstantTimeCompare, bukan ==: perbandingan
// biasa berhenti di byte pertama yang berbeda, dan selisih waktunya bisa
// dipakai menebak hash byte demi byte.
func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2-sha256" {
		return false, ErrInvalidHash
	}
	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations <= 0 {
		return false, ErrInvalidHash
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, ErrInvalidHash
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, ErrInvalidHash
	}

	got, err := pbkdf2.Key(sha256.New, password, salt, iterations, len(want))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
