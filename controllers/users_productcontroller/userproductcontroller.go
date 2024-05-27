package usersproductcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateUserProduct(w http.ResponseWriter, r *http.Request) {
	var UproductInput models.UsersProduct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&UproductInput); err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&UproductInput).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Memuat data pengguna dan produk yang terkait
	if err := models.DB.Preload("User").Preload("Product").First(&UproductInput, "id = ?", UproductInput.ID).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Mengembalikan respons dengan data yang dimuat
	response := map[string]interface{}{
		"message": "Data berhasil disimpan",
		"status":  true,
		"data":    UproductInput,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func ReadUserProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid, exists := params["uid"]

	// Jika uid tidak ada di URL, tampilkan seluruh data UsersProduct
	if !exists {
		var userProducts []models.UsersProduct
		if err := models.DB.Preload("User").Preload("Product").Find(&userProducts).Error; err != nil {
			response := map[string]interface{}{"message": err.Error(), "status": false}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data UsersProduct dalam format JSON
		if len(userProducts) > 0 {
			response := map[string]interface{}{"message": "Berhasil menampilkan data user_product", "status": true, "data": userProducts}
			helper.ResponseJSON(w, http.StatusOK, response)
			return
		} else {
			response := map[string]interface{}{"message": "Gagal menampilkan data user_product", "status": false, "data": userProducts}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Jika uid ada di URL, artinya kita ingin mengambil satu UserProduct berdasarkan ID
	// Konversi uidParam menjadi tipe data UUID
	uuid, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]interface{}{"message": "UID tidak valid", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data UserProduct berdasarkan ID
	var userProduct models.UsersProduct
	if err := models.DB.Preload("User").Preload("Product").Where("id = ?", uuid).First(&userProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]interface{}{"message": "User tidak ditemukan", "status": false}
			helper.ResponseJSON(w, http.StatusNotFound, response)
			return
		}
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Mengembalikan data User dalam format JSON
	response := map[string]interface{}{
		"message": "berhasil menampilkan data user_product",
		"status": true,
		"data":   userProduct,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func UpdateUserProduct(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter uid dari bagian path URL
	params := mux.Vars(r)
	uid := params["uid"]

	// Konversi uidParam menjadi tipe data UUID
	_, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]interface{}{"message": "UID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data userProduct yang akan diupdate
	var inputUserProduct models.UsersProduct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputUserProduct); err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// cek data userProduct berdasarkan UID
	existingUserProduct := models.UsersProduct{}
	if err := models.DB.First(&existingUserProduct, "id = ?", uid).Error; err != nil {
		response := map[string]interface{}{"message": "UID tidak ditemukan", "status": false}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Update data userProduct berdasarkan UID
	if err := models.DB.Model(&models.UsersProduct{}).Where("id = ?", uid).Updates(&inputUserProduct).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "userProduct berhasil diupdate", "status": true, "data": existingUserProduct}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteUserProduct(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter uid dari path URL
	vars := mux.Vars(r)
	uid := vars["uid"]

	// Konversi uidParam menjadi tipe data UUID
	id, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]interface{}{"message": "ID tidak valid", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Cek apakah data UsersProduct dengan ID tersebut ada
	var existingUsersProduct models.UsersProduct
	if err := models.DB.First(&existingUsersProduct, id).Error; err != nil {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan userUsersProduct dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.UsersProduct{}).Error; err != nil {
		response := map[string]interface{}{"message": errors.Wrap(err, "gagal menghapus User Product").Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "User Product berhasil dihapus", "status": true}
	helper.ResponseJSON(w, http.StatusOK, response)
}
