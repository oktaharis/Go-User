package usersrolecontroller

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

func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var URoleInput models.UsersRole
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&URoleInput); err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&URoleInput).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	data := []models.UsersRole{URoleInput}
	response := map[string]interface{}{
		"message": "success",
		"status" :  true,
		"data"	 :  data}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func ReaduserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")
	// Jika idParam tidak kosong, artinya kita ingin mengambil satu User berdasarkan ID
	if idParam != "" {
		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]interface{}{"message": "ID tidak valid", "status": false}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data userRole berdasarkan ID
		var userRole models.UsersRole
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Role").First(&userRole, id).Error; err != nil {
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
		helper.ResponseJSON(w, http.StatusOK, userRole)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data Users
		var userRoles []models.UsersRole
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Role").Find(&userRoles).Error; err != nil {
			response := map[string]interface{}{"message": err.Error(), "status": false}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data Users dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, map[string]interface{}{
			"status": true,
			 "data": userRoles})

	}
}
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter uid dari bagian path URL
	params := mux.Vars(r)
	uid := params["uid"]

	// Konversi uidParam menjadi tipe data UUID
	_, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]string{"message": "UID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data UserRole yang akan diupdate
	var inputUserRole models.UsersRole
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputUserRole); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Mencari data UserRole berdasarkan UID
	existingUserRole := models.UsersRole{}
	if err := models.DB.First(&existingUserRole, "id = ?", uid).Error; err != nil {
		response := map[string]string{"message": "UID tidak ditemukan"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Update data UserRole berdasarkan UID
	if err := models.DB.Model(&models.UsersRole{}).Where("id = ?", uid).Updates(&inputUserRole).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "UserRole berhasil diupdate", "status": true}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// Di dalam fungsi DeleteuserRole
func DeleteuserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter uid dari path URL
	vars := mux.Vars(r)
	uid := vars["uid"]

	// Konversi uidParam menjadi tipe data UUID
	id, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Cek apakah data UsersRole dengan ID tersebut ada
	var existingUsersRole models.UsersRole
	if err := models.DB.First(&existingUsersRole, id).Error; err != nil {
		response := map[string]string{"message": "Failed, ID tidak ditemukan"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan userUsersRole dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.UsersRole{}).Error; err != nil {
		response := map[string]string{"message": errors.Wrap(err, "gagal menghapus userRole").Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "userRole berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}