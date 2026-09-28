package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-food-store/auth"
	"go-food-store/helpers"
	"go-food-store/json/response"
)

type contextKey struct{}

// claimsKey memakai tipe privat sebagai kunci context, bukan string: kunci
// string bisa bentrok diam-diam dengan paket lain yang memakai teks yang sama.
var claimsKey = contextKey{}

// ClaimsFrom mengambil identitas pemanggil dari context. found bernilai false
// kalau request tidak melewati RequireAuth.
func ClaimsFrom(ctx context.Context) (auth.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(auth.Claims)
	return claims, ok
}

func denied(w http.ResponseWriter, code int, status, message string) {
	helpers.WriteJSON(w, response.WebResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    nil,
	})
}

// RequireAuth menolak request tanpa token yang sah dan menitipkan identitas
// pemanggil ke context.
func RequireAuth(issuer *auth.Issuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, found := strings.CutPrefix(header, "Bearer ")
			if !found || strings.TrimSpace(token) == "" {
				denied(w, http.StatusUnauthorized, "Unauthorized", "Missing bearer token")
				return
			}

			claims, err := issuer.Verify(strings.TrimSpace(token))
			if errors.Is(err, auth.ErrExpiredToken) {
				// Dibedakan dari token palsu: client yang tahu token-nya
				// kedaluwarsa bisa langsung login ulang tanpa menebak.
				denied(w, http.StatusUnauthorized, "Unauthorized", "Token expired, please login again")
				return
			} else if err != nil {
				denied(w, http.StatusUnauthorized, "Unauthorized", "Invalid token")
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey, claims)))
		})
	}
}

// RequireRole membatasi route ke peran tertentu. Dipasang setelah RequireAuth;
// tanpa itu tidak ada claims di context dan semua request ditolak 403.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFrom(r.Context())
			if !ok {
				denied(w, http.StatusForbidden, "Forbidden", "Not allowed to access this resource")
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			denied(w, http.StatusForbidden, "Forbidden", "Not allowed to access this resource")
		})
	}
}
