package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go-food-store/helpers"
	"go-food-store/json/response"
)

// Recover menangkap panic dari handler mana pun dan membalasnya sebagai 500
// berbentuk JSON.
//
// Dibutuhkan karena helpers.PanicHelper memang panic di tengah handler.
// net/http sendiri memang memulihkan panic per koneksi, tapi ia melakukannya
// dengan memutus koneksi begitu saja — client melihat "connection reset",
// bukan response. Di sini panic-nya dicatat sekali dan client tetap menerima
// envelope yang sama seperti error lain.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic pada handler",
					"method", r.Method,
					"path", r.URL.Path,
					"panic", rec,
				)
				helpers.WriteJSON(w, response.WebResponse{
					Code:    http.StatusInternalServerError,
					Status:  "Internal Server Error",
					Message: "Something went wrong",
					Data:    nil,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder mengintip status code yang ditulis handler.
//
// http.ResponseWriter tidak menyediakan cara membacanya kembali, jadi log
// akses tidak punya sumber lain. Kalau handler tidak pernah memanggil
// WriteHeader, net/http mengirim 200 — itulah nilai awalnya di sini.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Log mencatat satu baris per request: method, path, status, dan durasi.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"durasi", time.Since(start).Round(time.Millisecond).String(),
		)
	})
}

// CORS mengizinkan frontend yang dilayani dari origin berbeda memanggil API
// ini, dan menjawab preflight sebelum request aslinya dikirim.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Preflight berhenti di sini: ia tidak boleh diteruskan ke
			// handler, dan router akan membalasnya 405 karena tidak ada
			// route yang mendaftarkan OPTIONS.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
