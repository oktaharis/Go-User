package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeypc/go-jwt-mux/config"
	"github.com/jeypc/go-jwt-mux/helper"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]interface{}{"message": "Unauthorized, Harap Login Terlebih Dahulu!", "status": false}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		// mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaim{}
		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// token invalid
				response := map[string]interface{}{"status": false, "message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ValidationErrorExpired:
				// token expired
				response := map[string]interface{}{"status": false, "message": "Unauthorized, Token expired!"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]interface{}{"status": false, "message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]interface{}{"status": false, "message": "Unauthorized"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// Jika semua berjalan lancar, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
func ReadUserWithJWT(w http.ResponseWriter, r *http.Request) {
	tokenString := extractJWTFromURL(r.URL.Path)
	if tokenString == "" {
		response := map[string]interface{}{"status": false, "message": "Jwt token not found in url"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
	}

	// Parse JWT token
	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})
	if err != nil {
		response := map[string]interface{}{"status": false, "message": "Failed to parse jwt token"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
	}

	if !token.Valid {
		response := map[string]interface{}{"status": false, "message": "Invalid jwt token"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
	}

	// At this point, JWT is valid, you can access user profile data from `claims` variable
	profileData, err := json.Marshal(claims)
	if err != nil {
		response := map[string]interface{}{"status": false, "message": "Failed to marshal profile data"+err.Error()}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(profileData)
}

func extractJWTFromURL(urlPath string) string {
	parts := strings.Split(urlPath, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}