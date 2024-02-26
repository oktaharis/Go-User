package middlewares

import (
	"net/http"

	"github.com/jeypc/go-jwt-mux/config"
	"github.com/jeypc/go-jwt-mux/helper"
)
func SuperAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*config.JWTClaim)
		
		// Cek apakah pengguna memiliki peran "Super Admin"
		if claims.Role != "Super Admin" {
			response := map[string]interface{}{"status": false, "message": "Unauthorized, Anda tidak memiliki izin untuk mengakses sumber daya ini!"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Jika semua berjalan lancar, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
