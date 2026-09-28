package routes

import (
	"net/http"

	"gorm.io/gorm"

	"go-food-store/helpers"
	"go-food-store/json/response"
)

// Health melaporkan apakah proses ini bisa melayani request — termasuk apakah
// database-nya masih menjawab.
//
// Sebelumnya healthcheck docker-compose menembak /food/, yang berarti setiap
// tiga detik ada query SELECT ke tabel menu hanya untuk menanyakan "kamu
// hidup?". Ping jauh lebih murah dan tidak bercampur dengan data bisnis.
func Health(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.PingContext(r.Context())
		}
		if err != nil {
			helpers.WriteJSON(w, response.WebResponse{
				Code:    http.StatusServiceUnavailable,
				Status:  "Service Unavailable",
				Message: "Database is not reachable",
				Data:    nil,
			})
			return
		}
		helpers.WriteJSON(w, response.WebResponse{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Healthy",
			Data:    nil,
		})
	}
}
