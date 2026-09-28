package auth

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestHashPasswordMenghasilkanHashBerbedaTiapKali(t *testing.T) {
	// Salt acak berarti dua pengguna dengan kata sandi sama tidak punya hash
	// yang sama — tanpa itu, satu tebakan benar membongkar banyak akun.
	a, err := HashPassword("rahasia")
	if err != nil {
		t.Fatalf("hash gagal: %v", err)
	}
	b, err := HashPassword("rahasia")
	if err != nil {
		t.Fatalf("hash gagal: %v", err)
	}
	if a == b {
		t.Error("dua hash untuk kata sandi sama ternyata identik, salt tidak acak")
	}
	if strings.Contains(a, "rahasia") {
		t.Error("kata sandi bocor ke dalam hash")
	}
}

func TestVerifyPassword(t *testing.T) {
	hash, err := HashPassword("kata-sandi-benar")
	if err != nil {
		t.Fatalf("hash gagal: %v", err)
	}

	cases := []struct {
		nama  string
		input string
		cocok bool
	}{
		{"kata sandi benar", "kata-sandi-benar", true},
		{"kata sandi salah", "kata-sandi-salah", false},
		{"kata sandi kosong", "", false},
		{"beda kapitalisasi", "Kata-Sandi-Benar", false},
	}

	for _, c := range cases {
		t.Run(c.nama, func(t *testing.T) {
			cocok, err := VerifyPassword(c.input, hash)
			if err != nil {
				t.Fatalf("verifikasi error: %v", err)
			}
			if cocok != c.cocok {
				t.Errorf("cocok = %v, mau %v", cocok, c.cocok)
			}
		})
	}
}

func TestVerifyPasswordMenolakHashRusak(t *testing.T) {
	for _, rusak := range []string{"", "bukan-hash", "pbkdf2-sha256$nol$x$y", "md5$1$a$b"} {
		if _, err := VerifyPassword("apa saja", rusak); !errors.Is(err, ErrInvalidHash) {
			t.Errorf("hash %q: error = %v, mau ErrInvalidHash", rusak, err)
		}
	}
}

func TestTokenBolakBalik(t *testing.T) {
	issuer := NewIssuer("kunci-rahasia", time.Hour)

	token, err := issuer.Issue(7, "kasir", "admin")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}

	claims, err := issuer.Verify(token)
	if err != nil {
		t.Fatalf("verifikasi gagal: %v", err)
	}
	if claims.UserId != 7 || claims.Username != "kasir" || claims.Role != "admin" {
		t.Errorf("claims = %+v, tidak sesuai yang diterbitkan", claims)
	}
}

// Inti dari seluruh mekanisme ini: isi token tidak boleh bisa diubah tanpa
// kuncinya.
func TestVerifyMenolakTokenYangDirusak(t *testing.T) {
	issuer := NewIssuer("kunci-rahasia", time.Hour)
	token, err := issuer.Issue(1, "budi", "customer")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}
	body, signature, _ := strings.Cut(token, ".")

	cases := map[string]string{
		"tanpa titik":            body + signature,
		"tanda tangan diganti":   body + ".aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaA",
		"tanda tangan dihapus":   body + ".",
		"payload diganti":        "eyJzdWIiOjk5OSwicm9sZSI6ImFkbWluIiwiZXhwIjo5OTk5OTk5OTk5fQ." + signature,
		"token kosong":           "",
		"diterbitkan kunci lain": mustIssue(t, NewIssuer("kunci-penyerang", time.Hour)),
	}

	for nama, rusak := range cases {
		t.Run(nama, func(t *testing.T) {
			if _, err := issuer.Verify(rusak); err == nil {
				t.Error("token rusak diterima sebagai sah")
			}
		})
	}
}

func TestVerifyMenolakTokenKedaluwarsa(t *testing.T) {
	// ttl negatif: token lahir sudah lewat masa berlakunya, jadi test ini
	// tidak perlu menunggu.
	issuer := NewIssuer("kunci-rahasia", -time.Minute)
	token, err := issuer.Issue(1, "budi", "customer")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}

	if _, err := issuer.Verify(token); !errors.Is(err, ErrExpiredToken) {
		t.Errorf("error = %v, mau ErrExpiredToken", err)
	}
}

func mustIssue(t *testing.T, issuer *Issuer) string {
	t.Helper()
	token, err := issuer.Issue(1, "budi", "customer")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}
	return token
}
