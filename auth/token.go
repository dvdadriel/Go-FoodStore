package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("token tidak sah")
	ErrExpiredToken = errors.New("token sudah kedaluwarsa")
)

// Claims adalah isi token: siapa, perannya apa, dan berlaku sampai kapan.
type Claims struct {
	UserId   uint   `json:"sub"`
	Username string `json:"usr"`
	Role     string `json:"role"`
	ExpireAt int64  `json:"exp"`
}

// Issuer menerbitkan dan memverifikasi token akses.
//
// ponytail: token bertanda tangan HMAC-SHA256 buatan sendiri, bukan JWT dari
// pustaka pihak ketiga. Satu algoritma, tidak ada header algoritma yang bisa
// dipalsukan, dan seluruhnya pustaka standar — sekitar lima puluh baris. Ganti
// ke JWT betulan kalau nanti ada layanan lain yang harus ikut memverifikasi
// token ini; selama hanya API ini yang membacanya, formatnya tidak perlu
// standar.
//
// Yang belum ada: pencabutan token. Token yang bocor tetap berlaku sampai
// kedaluwarsa. Tambahkan daftar cabut kalau masa berlakunya dipanjangkan dari
// yang sekarang.
type Issuer struct {
	secret []byte
	ttl    time.Duration
}

func NewIssuer(secret string, ttl time.Duration) *Issuer {
	return &Issuer{secret: []byte(secret), ttl: ttl}
}

func (i *Issuer) sign(payload []byte) string {
	mac := hmac.New(sha256.New, i.secret)
	mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// Issue menerbitkan token untuk seorang pengguna.
func (i *Issuer) Issue(userId uint, username, role string) (string, error) {
	payload, err := json.Marshal(Claims{
		UserId:   userId,
		Username: username,
		Role:     role,
		ExpireAt: time.Now().Add(i.ttl).Unix(),
	})
	if err != nil {
		return "", err
	}
	body := base64.RawURLEncoding.EncodeToString(payload)
	return body + "." + i.sign(payload), nil
}

// Verify memeriksa tanda tangan dan masa berlaku, lalu mengembalikan isinya.
func (i *Issuer) Verify(token string) (Claims, error) {
	body, signature, found := strings.Cut(token, ".")
	if !found {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	// Tanda tangan diperiksa sebelum payload-nya dibaca: isi yang belum
	// terbukti asli tidak boleh mempengaruhi apa pun, termasuk parser JSON.
	if !hmac.Equal([]byte(signature), []byte(i.sign(payload))) {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if time.Now().Unix() >= claims.ExpireAt {
		return Claims{}, ErrExpiredToken
	}
	return claims, nil
}
