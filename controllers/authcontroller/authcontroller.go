package authcontroller

import (
	"crypto/rand"
	"database/sql"
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

	// Jika password belum di-hash
	if user.Password == userInput.Password {
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
		return
	}

	// Password sudah di-hash, lakukan validasi menggunakan bcrypt
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

// Fungsi untuk menambahkan token ke database
func addTokenToDatabase(name, token string, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO access_service (name, token) VALUES ('%s', '%s')", name, token)
	_, err := db.Exec(query)
	return err
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
	// Simpan name dan token ke dalam tabel access_service
	dbInstance, err := models.DB.DB()
	// fmt.Println("ini adalah dbinstance",dbInstance)
	if err != nil {
		response := map[string]string{"status": "failed", "message": "Gagal mendapatkan koneksi database"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	if err := addTokenToDatabase(user.Name, token, dbInstance); err != nil {
		fmt.Println(user.Name, token)
		response := map[string]string{"status": "failed", "message": "Gagal menyimpan name dan token ke database"}
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
	var product models.Product
	if err := models.DB.First(&role, user.RoleID).Error; err != nil {
		return "", err
	}
	// Query product data
	if err := models.DB.First(&product, user.ProductID).Error; err != nil {
		return "", err
	}
	// Set nilai Role dan Name dari pengguna sesuai dengan data dari database
	userRole := role.Name
	userProduct := product.Name
	userName := user.Name

	// Buat token JWT
	expTime := time.Now().Add(time.Minute * 72)
	claims := &config.JWTClaim{
		Email:   user.Email,
		Role:    userRole,
		Name:    userName,
		Product: userProduct,
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
	// Mengambil id dari path parameter
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Jika tidak ada id yang diberikan dalam URL, maka akan mengembalikan keseluruhan data user
		var users []models.User
		if err := models.DB.Preload("Role").Preload("Product").Find(&users).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		for i := range users {
			users[i].Product.ID = users[i].ProductID
		}

		response := map[string]interface{}{"message": "Menampilkan Seluruh data", "status": true, "data": users}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}

	// Jika id ditemukan dalam URL, maka akan mencari user dengan id tersebut dan mengembalikan data user tersebut
	var user models.User
	if err := models.DB.Preload("Role").Preload("Product").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"message": "User tidak ditemukan"}
			helper.ResponseJSON(w, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	// user.Product.ID = user.ProductID
	response := map[string]interface{}{"status": true, "data": user}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Mengambil input JSON yang diterima dari Postman
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Periksa apakah pengguna sudah ada di database berdasarkan alamat email
	var existingUser models.User
	if err := models.DB.Where("email = ?", userInput.Email).Preload("Role").First(&existingUser).Error; err == nil {
		// Jika pengguna sudah ada, kirim respons bahwa pengguna sudah dibuat
		response := map[string]interface{}{
			"message": "Pengguna sudah dibuat",
			"status":  false,
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}

	// Hash password menggunakan bcrypt sebelum menyimpannya ke database
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

	// Mengambil data pengguna yang baru dibuat dari database bersama dengan relasi yang sesuai
	var newUser models.User
	if err := models.DB.Preload("Role").Preload("Product").First(&newUser, userInput.ID).Error; err != nil {
		response := map[string]string{"message": "Gagal mengambil data pengguna yang baru dibuat"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{
		"message": "Pengguna berhasil terdaftar",
		"status":  true,
		"data":    newUser,
	}
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
	// Mendapatkan parameter id dari bagian path URL
	params := mux.Vars(r)
	idParam, ok := params["id"]
	if !ok {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan dalam path URL", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]interface{}{"message": "Failed, ID tidak valid", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Cek apakah data Product dengan ID tersebut ada
	var existingUser models.User
	if err := models.DB.Preload("Role").Preload("Product").First(&existingUser, id).Error; err != nil {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data user yang akan diupdate
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]interface{}{"message": "Gagal menguraikan payload JSON"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Hashing password baru sebelum menyimpannya ke database
	if userInput.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			response := map[string]interface{}{"message": "Gagal menghash password baru"}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
		userInput.Password = string(hashedPassword)
	}

	// Update data user berdasarkan ID
	if err := models.DB.Model(&models.User{}).Where("id = ?", id).Updates(&userInput).Error; err != nil {
		response := map[string]interface{}{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "User berhasil diupdate", "status": true, "data": existingUser}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari path variabel
	vars := mux.Vars(r)
	idParam := vars["id"]

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]interface{}{"message": "ID tidak valid", "status": false}
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
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
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
			response := map[string]string{"status": "failed", "message": "Email tidak ditemukan. Pastikan email yang Anda masukkan benar."}
			helper.ResponseJSON(w, http.StatusNotFound, response)
		default:
			response := map[string]string{"status": "failed", "message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
		}
		return
	}

	// Generate password baru
	newPassword := generateNewPassword(user.Name, user.Email)

	// Simpan password baru ke dalam database
	user.Password = newPassword
	if err := models.DB.Save(&user).Error; err != nil {
		response := map[string]string{"status": "failed", "message": "Gagal menyimpan password baru"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Kirim email ke pengguna dengan password baru
	// Implementasi logika pengiriman email Anda di sini

	response := map[string]interface{}{
		"status":  "success",
		"message": "Password baru telah berhasil dibuat dan dikirim ke email Anda.",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func generateNewPassword(username, email string) string {
	// Ambil 3 karakter pertama dari nama pengguna dan email
	usernamePrefix := username[:3]
	emailPrefix := email[:3]

	// Gabungkan 3 karakter pertama dari nama pengguna dan email untuk membuat password baru
	newPassword := usernamePrefix + emailPrefix

	// Jika panjang password kurang dari 6 karakter, tambahkan karakter acak hingga panjangnya menjadi 6
	for len(newPassword) < 6 {
		newPassword += string(randomChar())
	}

	return newPassword
}

func randomChar() byte {
	// Karakter yang diperbolehkan untuk digunakan dalam password
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Mendapatkan indeks karakter acak menggunakan crypto/rand
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	return chars[idx.Int64()]
}
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Dapatkan ID pengguna dari parameter query HTTP
	params := mux.Vars(r)
	idParam, ok := params["id"]

	if !ok {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan dalam path URL", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi idParam ke tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]interface{}{"status": "error", "message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Dapatkan data pengguna yang akan di-reset berdasarkan ID
	var user models.User
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		response := map[string]interface{}{"status": "error", "message": "Pengguna tidak ditemukan"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Parse token reset dan password baru dari body request JSON
	var resetInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&resetInput); err != nil {
		response := map[string]interface{}{"status": "error", "message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Memperbarui data pengguna di database berdasarkan ID dengan password baru
	if resetInput.ResetPassword != "" {
		// Jika password baru ada, gunakan password baru tanpa meng-hash
		user.Password = resetInput.ResetPassword
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", id).Updates(&user).Error; err != nil {
		response := map[string]interface{}{"status": "error", "message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Memberikan respons JSON berisi pesan bahwa password telah berhasil di-reset
	response := map[string]interface{}{"status": "success", "message": "Password berhasil direset"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
