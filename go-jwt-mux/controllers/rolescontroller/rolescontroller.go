package rolescontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var roleInput models.Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&roleInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&roleInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func GetRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Jika idParam tidak kosong, artinya kita ingin mengambil satu role berdasarkan ID
	if idParam != "" {
		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data role berdasarkan ID
		var role models.Role
		if err := models.DB.First(&role, id).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan data role dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, role)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data roles
		var roles []models.Role
		if err := models.DB.Find(&roles).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data roles dalam format JSON
		helper.ResponseJSON(w, http.StatusOK, roles)
	}
}


func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data role yang akan diupdate
	var roleInput models.Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&roleInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	// Update data Role berdasarkan ID
	if err := models.DB.Model(&models.Role{}).Where("id = ?", id).Updates(&roleInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Role berhasil diupdate"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan Role dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.Role{}).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Role berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
