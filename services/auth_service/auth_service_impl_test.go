package authservice

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"

	"go-food-store/auth"
	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
)

const testSecret = "kunci-test-yang-panjangnya-cukup-32-karakter"

// stubUserRepo dikendalikan test; field fungsi nil berarti method-nya tidak
// boleh dipanggil.
type stubUserRepo struct {
	t *testing.T

	createFn func(models.User) (models.User, error)
	findFn   func(string) (models.User, error)
	existsFn func(string) (bool, error)

	created models.User
	calls   []string
}

func (s *stubUserRepo) Create(user models.User) (models.User, error) {
	s.t.Helper()
	if s.createFn == nil {
		s.t.Fatalf("Create dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "Create")
	s.created = user
	return s.createFn(user)
}

func (s *stubUserRepo) FindByUsername(username string) (models.User, error) {
	s.t.Helper()
	if s.findFn == nil {
		s.t.Fatalf("FindByUsername dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "FindByUsername")
	return s.findFn(username)
}

func (s *stubUserRepo) ExistsByUsername(username string) (bool, error) {
	s.t.Helper()
	if s.existsFn == nil {
		s.t.Fatalf("ExistsByUsername dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "ExistsByUsername")
	return s.existsFn(username)
}

func (s *stubUserRepo) called(name string) bool {
	for _, c := range s.calls {
		if c == name {
			return true
		}
	}
	return false
}

func newService(t *testing.T, repo *stubUserRepo) AuthService {
	t.Helper()
	return NewAuthService(repo, auth.NewIssuer(testSecret, time.Hour), 3600, validator.New())
}

// Pendaftaran terbuka tidak boleh bisa menghasilkan admin, apa pun isi
// request-nya. Ini satu-satunya hal yang menahan siapa pun jadi admin.
func TestRegisterSelaluMenghasilkanPeranCustomer(t *testing.T) {
	repo := &stubUserRepo{t: t}
	repo.existsFn = func(string) (bool, error) { return false, nil }
	repo.createFn = func(u models.User) (models.User, error) {
		u.Model = gorm.Model{ID: 1}
		return u, nil
	}

	res := newService(t, repo).Register(request.RegisterReq{
		Username: "budi",
		Password: "kata-sandi-panjang",
	})

	if res.Code != http.StatusCreated {
		t.Fatalf("Code = %d, mau %d", res.Code, http.StatusCreated)
	}
	if repo.created.Role != models.RoleCustomer {
		t.Errorf("peran tersimpan = %q, mau %q", repo.created.Role, models.RoleCustomer)
	}
	if repo.created.PasswordHash == "kata-sandi-panjang" || repo.created.PasswordHash == "" {
		t.Error("kata sandi tidak di-hash sebelum disimpan")
	}

	user, ok := res.Data.(response.UserResponse)
	if !ok {
		t.Fatalf("Data bertipe %T, mau response.UserResponse", res.Data)
	}
	if user.Role != models.RoleCustomer {
		t.Errorf("peran di response = %q, mau %q", user.Role, models.RoleCustomer)
	}
}

func TestRegisterMenolakInputTidakValid(t *testing.T) {
	cases := map[string]request.RegisterReq{
		"username terlalu pendek": {Username: "ab", Password: "kata-sandi-panjang"},
		"password terlalu pendek": {Username: "budi", Password: "pendek"},
		"keduanya kosong":         {},
	}

	for nama, req := range cases {
		t.Run(nama, func(t *testing.T) {
			repo := &stubUserRepo{t: t}
			res := newService(t, repo).Register(req)

			if res.Code != http.StatusBadRequest {
				t.Errorf("Code = %d, mau %d", res.Code, http.StatusBadRequest)
			}
			if repo.called("Create") {
				t.Error("user tetap dibuat padahal input tidak valid")
			}
		})
	}
}

func TestRegisterMenolakUsernameTerpakai(t *testing.T) {
	repo := &stubUserRepo{t: t}
	repo.existsFn = func(string) (bool, error) { return true, nil }

	res := newService(t, repo).Register(request.RegisterReq{
		Username: "budi",
		Password: "kata-sandi-panjang",
	})

	if res.Code != http.StatusConflict {
		t.Errorf("Code = %d, mau %d", res.Code, http.StatusConflict)
	}
	if repo.called("Create") {
		t.Error("user tetap dibuat padahal username sudah terpakai")
	}
}

func TestLoginMengembalikanTokenYangBisaDiverifikasi(t *testing.T) {
	hash, err := auth.HashPassword("kata-sandi-panjang")
	if err != nil {
		t.Fatalf("hash gagal: %v", err)
	}
	repo := &stubUserRepo{t: t}
	repo.findFn = func(string) (models.User, error) {
		return models.User{
			Model:        gorm.Model{ID: 9},
			Username:     "kasir",
			PasswordHash: hash,
			Role:         models.RoleAdmin,
		}, nil
	}

	res := newService(t, repo).Login(request.LoginReq{Username: "kasir", Password: "kata-sandi-panjang"})
	if res.Code != http.StatusOK {
		t.Fatalf("Code = %d, mau %d", res.Code, http.StatusOK)
	}

	login, ok := res.Data.(response.LoginResponse)
	if !ok {
		t.Fatalf("Data bertipe %T, mau response.LoginResponse", res.Data)
	}

	claims, err := auth.NewIssuer(testSecret, time.Hour).Verify(login.Token)
	if err != nil {
		t.Fatalf("token yang diterbitkan tidak lolos verifikasi: %v", err)
	}
	if claims.UserId != 9 || claims.Role != models.RoleAdmin {
		t.Errorf("claims = %+v, tidak sesuai user yang login", claims)
	}
}

// Username tidak terdaftar dan kata sandi salah harus tidak terbedakan dari
// luar. Pesan yang berbeda memberi tahu penyerang username mana yang ada.
func TestLoginTidakMembocorkanUsernameManaYangAda(t *testing.T) {
	hash, err := auth.HashPassword("kata-sandi-benar")
	if err != nil {
		t.Fatalf("hash gagal: %v", err)
	}

	userTidakAda := &stubUserRepo{t: t}
	userTidakAda.findFn = func(string) (models.User, error) {
		return models.User{}, gorm.ErrRecordNotFound
	}

	sandiSalah := &stubUserRepo{t: t}
	sandiSalah.findFn = func(string) (models.User, error) {
		return models.User{Model: gorm.Model{ID: 1}, Username: "budi", PasswordHash: hash, Role: models.RoleCustomer}, nil
	}

	a := newService(t, userTidakAda).Login(request.LoginReq{Username: "hantu", Password: "apa-saja"})
	b := newService(t, sandiSalah).Login(request.LoginReq{Username: "budi", Password: "kata-sandi-salah"})

	if a.Code != http.StatusUnauthorized || b.Code != http.StatusUnauthorized {
		t.Fatalf("Code = %d dan %d, keduanya mau %d", a.Code, b.Code, http.StatusUnauthorized)
	}
	if a.Message != b.Message {
		t.Errorf("pesan berbeda (%q vs %q) — membocorkan username mana yang terdaftar", a.Message, b.Message)
	}
	if a.Data != nil || b.Data != nil {
		t.Error("login gagal tetap mengembalikan Data")
	}
}

func TestLoginMelaporkanErrorDatabaseApaAdanya(t *testing.T) {
	repo := &stubUserRepo{t: t}
	repo.findFn = func(string) (models.User, error) {
		return models.User{}, errors.New("koneksi putus")
	}

	res := newService(t, repo).Login(request.LoginReq{Username: "budi", Password: "apa-saja"})

	// 500, bukan 401: database yang mati bukan berarti kredensialnya salah.
	if res.Code != http.StatusInternalServerError {
		t.Errorf("Code = %d, mau %d", res.Code, http.StatusInternalServerError)
	}
}
