package usersrolecontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var URoleInput models.UsersRole
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&URoleInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&URoleInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	response := map[string]string{"message": "success"}
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
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data userRole berdasarkan ID
		var userRole models.UsersRole
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Role").First(&userRole, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response := map[string]string{"message": "User tidak ditemukan"}
				helper.ResponseJSON(w, http.StatusNotFound, response)
				return
			}
			response := map[string]string{"message": err.Error()}
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
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data Users dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, userRoles)
	}
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data UUID
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data userRole yang akan diupdate
	var inputuserRole models.UsersRole
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputuserRole); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Mencari data userRole berdasarkan ID
	existinguserRole := models.UsersRole{}
	if err := models.DB.First(&existinguserRole, "id = ?", id).Error; err != nil {
		response := map[string]string{"message": "ID tidak ditemukan"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Update data userRole berdasarkan ID
	if err := models.DB.Model(&models.UsersRole{}).Where("id = ?", id).Updates(&inputuserRole).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "userRole berhasil diupdate"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteuserRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data UUID
	id, err := uuid.Parse(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan userRole dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.UsersRole{}).Error; err != nil {
		response := map[string]string{"message": errors.Wrap(err, "gagal menghapus userRole").Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "userRole berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
