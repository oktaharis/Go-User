package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeypc/go-jwt-mux/config"
	"github.com/jeypc/go-jwt-mux/helper"
)

func  RoleAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil informasi token dari cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// Jika tidak ada cookie "token" yang ditemukan, maka pengguna tidak terotentikasi
			response := map[string]interface{}{"message": "Unauthorized, cookie tidak di temukan! ", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Parse token JWT dari cookie
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})
		if err != nil {
			// Jika terjadi kesalahan dalam parsing token, pengguna tidak terotentikasi
			response := map[string]interface{}{"message": "Unauthorized, Pengguna tidak terotentikasi!", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Verifikasi token JWT
		if !token.Valid {
			// Jika token tidak valid, pengguna tidak terotentikasi
			response := map[string]interface{}{"message": "Unauthorized, Token tidak valid!", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Ambil klaim (claims) dari token JWT
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// Jika tidak dapat membaca klaim dari token, pengguna tidak terotentikasi
			response := map[string]interface{}{"message": "Unauthorized, tidak dapat membaca klaim dari token!", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Periksa apakah peran pengguna adalah Super Admin
		role, ok := claims["Role"].(string)
		if !ok {
			// Jika tidak dapat membaca peran dari klaim, pengguna tidak terotentikasi
			response := map[string]interface{}{"message": "Unauthorized, Role Kosong!", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Jika peran bukan Super Admin, tolak akses
		if role != "Super Admin" {
			response := map[string]interface{}{"message": "Unauthorized, Harap Login Terlebih Dahulu!", "status": false}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Jika peran adalah Super Admin, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
