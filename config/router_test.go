package config_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"

	"go-food-store/config"
)

// Test tingkat router: membuktikan bahwa status HTTP yang benar benar-benar
// sampai ke kabel, lewat rantai middleware yang sebenarnya — bukan hanya
// benar di dalam helper.
//
// db sengaja nil. Route yang diuji di sini tidak menyentuh database sama
// sekali; yang diuji adalah routing dan middleware-nya.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
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
