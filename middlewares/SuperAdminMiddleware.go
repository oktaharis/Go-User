package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeypc/go-jwt-mux/config"
)

func RoleAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil informasi token dari cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// Jika tidak ada cookie "token" yang ditemukan, maka pengguna tidak terotentikasi
			http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
			return
		}

		// Parse token JWT dari cookie
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})
		if err != nil {
			// Jika terjadi kesalahan dalam parsing token, pengguna tidak terotentikasi
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Verifikasi token JWT
		if !token.Valid {
			// Jika token tidak valid, pengguna tidak terotentikasi
			http.Error(w, "Unauthorized: Token is not valid", http.StatusUnauthorized)
			return
		}

		// Ambil klaim (claims) dari token JWT
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// Jika tidak dapat membaca klaim dari token, pengguna tidak terotentikasi
			http.Error(w, "Unauthorized: Unable to read token claims", http.StatusUnauthorized)
			return
		}

		// Periksa apakah peran pengguna adalah Super Admin
		role, ok := claims["Role"].(string)
		if !ok {
			// Jika tidak dapat membaca peran dari klaim, pengguna tidak terotentikasi
			http.Error(w, "Unauthorized: Role kosong", http.StatusUnauthorized)
			return
		}

		// Jika peran bukan Super Admin, tolak akses
		if role != "Super Admin" {
			http.Error(w, "Forbidden: Only Super Admin can perform this operation", http.StatusForbidden)
			return
		}

		// Jika peran adalah Super Admin, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

