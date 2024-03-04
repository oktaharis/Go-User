package authcontroller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

// Update fungsi Login untuk menggunakan GenerateJWT
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

	// Generate OTP
	otp := generateOTP()

	// Simpan OTP ke dalam database
	user.OTP = otp
	if err := models.DB.Save(&user).Error; err != nil {
		response := map[string]string{"status": "failed", "message": "Gagal menyimpan OTP"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Response ke pengguna bahwa OTP berhasil dibuat
	response := map[string]interface{}{
		"status":  "success",
		"message": "OTP berhasil dibuat",
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
			// Handle the error, for example, log it or return a default value
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

	// Jika OTP valid, buat token JWT menggunakan GenerateJWT
	token, err := GenerateJWT(user)
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

	// Response JSON dengan token JWT
	response := map[string]interface{}{
		"status":  "success",
		"message": "Login berhasil",
		"token":   token,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// Buat fungsi baru untuk membuat token JWT
func GenerateJWT(user models.User) (string, error) {
	// Memuat informasi peran dan nama pengguna dari database
	var role models.Role
	if err := models.DB.First(&role, user.RoleID).Error; err != nil {
		return "", err
	}

	// Set nilai Role dan Name dari pengguna sesuai dengan data dari database
	userRole := role.Name
	userName := user.Name

	// Buat token JWT
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Email: user.Email,
		Role:  userRole,
		Name:  userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// Buat token JWT
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		return "", err
	}

	return token, nil
}

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
	data := []models.User{userInput}
	response := map[string]interface{}{
		"message": "success",
		"status" :  true,
		"data"	 :  data}
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
	// Mengambil ID dari path variabel
	vars := mux.Vars(r)
	idParam := vars["id"]

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Cek apakah data User dengan ID tersebut ada
	var existingUser models.User
	if err := models.DB.First(&existingUser, id).Error; err != nil {
		response := map[string]string{"message": "Failed, ID tidak ditemukan"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan User dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "User berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}