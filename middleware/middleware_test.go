package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Panic di dalam handler harus jadi 500 berbentuk JSON, bukan koneksi putus.
func TestRecoverMembalas500BukanMemutusKoneksi(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("sesuatu meledak")
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/food/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, mau %d", rec.Code, http.StatusInternalServerError)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, mau application/json", ct)
	}
}

func TestRecoverTidakMengubahResponseNormal(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/food/", nil))

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, mau %d", rec.Code, http.StatusCreated)
	}
}

// Preflight harus dijawab middleware dan berhenti di situ. Kalau diteruskan,
// router membalas 405 karena tidak ada route yang mendaftarkan OPTIONS.
func TestCORSMenjawabPreflightTanpaMeneruskan(t *testing.T) {
	diteruskan := false
	handler := CORS("*")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		diteruskan = true
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodOptions, "/food/", nil))

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, mau %d", rec.Code, http.StatusNoContent)
	}
	if diteruskan {
		t.Error("preflight diteruskan ke handler")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, mau *", got)
	}
}

func TestCORSMeneruskanRequestBiasa(t *testing.T) {
	diteruskan := false
	handler := CORS("https://toko.example")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		diteruskan = true
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/food/", nil))

	if !diteruskan {
		t.Error("request biasa tidak diteruskan")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://toko.example" {
		t.Errorf("Allow-Origin = %q, mau https://toko.example", got)
	}
}
