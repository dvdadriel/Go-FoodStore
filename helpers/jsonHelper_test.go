package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-food-store/json/response"
)

// Regresi: sebelum perbaikan ini setiap handler membalas 200 apa pun isi
// envelope-nya, karena tidak ada yang memanggil WriteHeader. Kode error hanya
// hidup di body, sehingga client maupun CI tidak bisa membedakan sukses dari
// gagal.
func TestWriteJSONMemakaiCodeSebagaiStatusHTTP(t *testing.T) {
	cases := []struct {
		nama string
		code int
	}{
		{"sukses", http.StatusOK},
		{"tidak ditemukan", http.StatusNotFound},
		{"permintaan salah", http.StatusBadRequest},
		{"bentrok", http.StatusConflict},
		{"error server", http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.nama, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteJSON(rec, response.WebResponse{Code: c.code, Status: c.nama})

			if rec.Code != c.code {
				t.Errorf("status HTTP = %d, mau %d", rec.Code, c.code)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, mau application/json", ct)
			}

			var body response.WebResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body bukan JSON yang sah: %v", err)
			}
			if body.Code != c.code {
				t.Errorf("Code di body = %d, mau %d", body.Code, c.code)
			}
		})
	}
}
