package usersproductcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	data := []models.UsersProduct{UproductInput}
	response := map[string]interface{}{
		"message": "success",
		"status" :  true,
		"data"	 :  data}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func ReadUserProduct(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")
	// Jika idParam tidak kosong, artinya kita ingin mengambil satu User berdasarkan ID
	if idParam != "" {
		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]interface{}{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data UserProduct berdasarkan ID
		var userProduct models.UsersProduct
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Product").First(&userProduct, id).Error; err != nil {
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
		helper.ResponseJSON(w, http.StatusOK, userProduct)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data Users
		var userProducts []models.UsersProduct
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Product").Find(&userProducts).Error; err != nil {
			response := map[string]interface{}{"message": err.Error(), "status": false}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data Users dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, userProducts)
	}
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

	// Mencari data userProduct berdasarkan UID
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

	response := map[string]interface{}{"message": "userProduct berhasil diupdate", "status": true}
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
