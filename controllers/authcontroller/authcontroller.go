package authcontroller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/jeypc/go-jwt-mux/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jeypc/go-jwt-mux/helper"
	"gorm.io/gorm"

	"strconv"

	"github.com/jeypc/go-jwt-mux/models"
	"golang.org/x/crypto/bcrypt"
)
// func sendOTPByEmail(toEmail, otp string) error {
//     from := mail.NewEmail("Okta Haris Sutanto", "oktaharis2008@gmail.com")
//     subject := "Kode OTP untuk Verifikasi"
//     to := mail.NewEmail("", toEmail) // Menggunakan alamat email tujuan yang diberikan
//     plainTextContent := fmt.Sprintf("Kode OTP Anda adalah: %s", otp)
//     htmlContent := fmt.Sprintf("<strong>Kode OTP Anda adalah:</strong> %s", otp)
//     message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
//     client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
//     response, err := client.Send(message)
//     if err != nil {
//         log.Println("Error sending OTP email:", err)
//         return err
//     }
//     if response.StatusCode >= 200 && response.StatusCode < 300 {
//         log.Println("Email OTP berhasil dikirim:", response.StatusCode)
//     } else {
//         log.Println("Email OTP gagal dikirim:", response.StatusCode)
//     }
//     return nil
// }

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"status": "failed", "message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"status": "failed", "message": "Email atau password salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
		default:
			response := map[string]string{"status": "failed", "message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"status": "failed", "message": "Email atau password salah"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	generatedOTP := generateOTP()
	user.OTP = generatedOTP
	models.DB.Save(&user)
	  // Kirim email OTP menggunakan SendGrid
	//   if err := sendOTPByEmail(userInput.Email, generatedOTP); err != nil {
    //     http.Error(w, "Gagal mengirimkan email OTP", http.StatusInternalServerError)
    //     return
    // }

	response := map[string]interface{}{
		"status":  "success",
		"message": "Login berhasil",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func generateOTP() string {
	// Generate a random 6-digit OTP using the crypto/rand package
	const letterBytes = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			// Handle the error, for gmail, log it or return a default value
			return "000000"
		}
		otp[i] = letterBytes[randomIndex.Int64()]
	}
	return string(otp)
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var inputOTP struct {
		OTP string `json:"otp"`
	}

	// Decode input JSON
	var userInput models.User
	fmt.Println(userInput)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputOTP); err != nil {
		response := map[string]string{"status": "failed", "message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Ambil data pengguna dari database berdasarkan otp
	var user models.User
	if err := models.DB.Where("otp = ?", inputOTP.OTP).First(&user).Error; err != nil {
		log.Printf("Error querying database: %v", err)
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"status": "failed", "message": "OTP tidak sesuai dengan yang di database"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"status": "failed", "message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}
	
	// Validasi OTP
	if inputOTP.OTP != user.OTP {
		response := map[string]string{"status": "failed", "message": "Kode OTP tidak sesuai"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Jika OTP valid, buat token JWT
	// Proses pembuatan token jwt
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	// Medeklarasikan algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"status": "failed", "message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Hapus atau atur ulang OTP setelah digunakan
	models.DB.Save(&user)
	// Set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]interface{}{
		"status":  "success",
		"message": "Login berhasil",
		"token":   token,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// handlers/user_handler.go


func ReadUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")
	// Jika idParam tidak kosong, artinya kita ingin mengambil satu User berdasarkan ID
	if idParam != "" {
		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data User berdasarkan ID
		var user models.User
		if err := models.DB.Preload("Role").Preload("Product").First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response := map[string]string{"message": "User tidak ditemukan"}
				helper.ResponseJSON(w, http.StatusNotFound, response)
				return
			}
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengisi nilai id pada objek product dengan nilai dari product_id
		user.Product.ID = user.ProductID

		// Mengembalikan data User dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, user)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data Users
		var users []models.User
		if err := models.DB.Preload("Role").Preload("Product").Find(&users).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengisi nilai id pada setiap objek product dengan nilai dari product_id
		for i := range users {
			users[i].Product.ID = users[i].ProductID
		}

		// Mengembalikan seluruh data Users dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, users)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan json yang di terima dari postman
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal menghash password")
	}
	userInput.Password = string(hashedPassword)

	// INSERT KE DATABASE
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func Logout(w http.ResponseWriter, r *http.Request) { // perubahan logout menjadi Logout
	// hapus toke yang ada di cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Berhasil"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data user yang akan diupdate
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Hashing password baru sebelum menyimpannya ke database
	if userInput.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			response := map[string]string{"message": "Gagal menghash password baru"}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
		userInput.Password = string(hashedPassword)
	}

	// Update data user berdasarkan ID
	if err := models.DB.Model(&models.User{}).Where("id = ?", id).Updates(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User berhasil diupdate"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan user dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
