package config_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"

	"go-food-store/auth"
	"go-food-store/config"
)

// Test tingkat router: membuktikan bahwa status HTTP yang benar benar-benar
// sampai ke kabel, lewat rantai middleware yang sebenarnya — bukan hanya
// benar di dalam helper.
//
// db sengaja nil. Route yang diuji di sini tidak menyentuh database sama
// sekali; yang diuji adalah routing dan middleware-nya.
const testSecret = "kunci-test-yang-panjangnya-cukup-32-karakter"

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
	server := httptest.NewServer(config.SetupModel(nil, validator.New()))
	t.Cleanup(server.Close)
	return server
}

func TestRouterMengirimStatusHTTPYangBenar(t *testing.T) {
	server := newTestServer(t)

	cases := []struct {
		nama   string
		method string
		path   string
		status int
	}{
		{"route tidak terdaftar", http.MethodGet, "/nope", http.StatusNotFound},
		{"method salah", http.MethodPatch, "/food/", http.StatusMethodNotAllowed},
		{"id bukan angka", http.MethodGet, "/food/bukan-angka", http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.nama, func(t *testing.T) {
			req, err := http.NewRequest(c.method, server.URL+c.path, nil)
			if err != nil {
				t.Fatalf("gagal membuat request: %v", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request gagal: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != c.status {
				t.Errorf("status = %d, mau %d", res.StatusCode, c.status)
			}
			if ct := res.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, mau application/json", ct)
			}
		})
	}
}

// Handler yang panic karena database nil tidak boleh memutus koneksi:
// client harus tetap menerima response 500 yang utuh.
func TestRouterMembalas500SaatHandlerPanic(t *testing.T) {
	server := newTestServer(t)

	res, err := http.Get(server.URL + "/food/")
	if err != nil {
		t.Fatalf("koneksi terputus, Recover tidak bekerja: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, mau %d", res.StatusCode, http.StatusInternalServerError)
	}
}

// Route yang dijaga harus menolak sebelum menyentuh handler — kalau tidak,
// request tanpa token akan lolos ke database dan panik karena db nil di sini,
// yang akan terlihat sebagai 500 alih-alih 401.
func TestRouteTerjagaMenolakTanpaToken(t *testing.T) {
	server := newTestServer(t)

	paths := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/food/"},
		{http.MethodDelete, "/food/1"},
		{http.MethodGet, "/cust/"},
		{http.MethodPost, "/transaction/"},
		{http.MethodPut, "/transaction/acc/1"},
		{http.MethodGet, "/transaction/unpaid/"},
	}

	for _, p := range paths {
		t.Run(p.method+" "+p.path, func(t *testing.T) {
			res := do(t, server, p.method, p.path, "")
			defer res.Body.Close()
			if res.StatusCode != http.StatusUnauthorized {
				t.Errorf("status = %d, mau %d", res.StatusCode, http.StatusUnauthorized)
			}
		})
	}
}

func TestRouteTerjagaMenolakTokenPalsu(t *testing.T) {
	server := newTestServer(t)

	penyerang := auth.NewIssuer("kunci-lain-yang-panjangnya-juga-32-karakter", time.Hour)
	token, err := penyerang.Issue(1, "penyusup", "admin")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}

	res := do(t, server, http.MethodPost, "/food/", token)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, mau %d", res.StatusCode, http.StatusUnauthorized)
	}
}

// Customer yang sah tetap tidak boleh mengelola menu atau menerima pembayaran.
func TestRouteAdminMenolakCustomer(t *testing.T) {
	server := newTestServer(t)

	issuer := auth.NewIssuer(testSecret, time.Hour)
	token, err := issuer.Issue(2, "budi", "customer")
	if err != nil {
		t.Fatalf("penerbitan gagal: %v", err)
	}

	for _, p := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/food/"},
		{http.MethodPut, "/transaction/acc/1"},
		{http.MethodGet, "/transaction/paid/"},
	} {
		t.Run(p.method+" "+p.path, func(t *testing.T) {
			res := do(t, server, p.method, p.path, token)
			defer res.Body.Close()
			if res.StatusCode != http.StatusForbidden {
				t.Errorf("status = %d, mau %d", res.StatusCode, http.StatusForbidden)
			}
		})
	}
}

// Etalase tetap terbuka: daftar menu bisa dibaca tanpa login. db nil di sini
// membuat handler-nya panik, jadi yang dibuktikan adalah ia benar-benar
// sampai ke handler (500) dan tidak dihadang penjaga (401).
func TestDaftarMenuTerbukaTanpaToken(t *testing.T) {
	server := newTestServer(t)

	res := do(t, server, http.MethodGet, "/food/", "")
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		t.Error("daftar menu ikut terjaga, padahal harus terbuka")
	}
}

func do(t *testing.T, server *httptest.Server, method, path, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, server.URL+path, nil)
	if err != nil {
		t.Fatalf("gagal membuat request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request gagal: %v", err)
	}
	return res
}
