package usersrolecontroller

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

	// Memuat data pengguna dan produk yang terkait
	if err := models.DB.Preload("User").Preload("Role").First(&URoleInput, "id = ?", URoleInput.ID).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Mengembalikan respons dengan data yang dimuat
	response := map[string]interface{}{
		"message": "Data berhasil disimpan",
		"status":  true,
		"data":    URoleInput,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func ReadUserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan nilai parameter "uid" dari bagian path URL menggunakan Gorilla Mux
	params := mux.Vars(r)
	uid, exists := params["uid"]

	// Jika uid tidak ada di URL, tampilkan seluruh data UsersRole
	if !exists {
		var userRole []models.UsersRole
		if err := models.DB.Preload("User").Preload("Role").Find(&userRole).Error; err != nil {
			response := map[string]interface{}{"message": err.Error(), "status": false}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data UsersRole dalam format JSON
		if len(userRole) > 0 {
			response := map[string]interface{}{"message": "Berhasil menampilkan data user_role", "status": true, "data": userRole}
			helper.ResponseJSON(w, http.StatusOK, response)
			return
		} else {
			response := map[string]interface{}{"message": "Gagal menampilkan data user_role", "status": false, "data": userRole}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
		
	}

	// Jika uid ada di URL, artinya kita ingin mengambil satu UserRole berdasarkan ID
	// Konversi uidParam menjadi tipe data UUID
	_, err := uuid.Parse(uid)
	if err != nil {
		response := map[string]interface{}{"message": "UID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data UserRole berdasarkan ID
	var userRole models.UsersRole
	if err := models.DB.Preload("User").Preload("Role").Where("id = ?", uid).First(&userRole).Error; err != nil {
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
	response := map[string]interface{}{"message": "berhasil menampilkan data user_product",
		"status": false,
		"data":   userRole}
	helper.ResponseJSON(w, http.StatusInternalServerError, response)
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
		response := map[string]interface{}{"message": "UID tidak ditemukan", "status": false}
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

	response := map[string]interface{}{"message": "userRole berhasil dihapus", "status": true}
	helper.ResponseJSON(w, http.StatusOK, response)
}
