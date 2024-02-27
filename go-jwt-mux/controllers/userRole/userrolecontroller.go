package userrolecontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
)

// CreateRole membuat role baru
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var roleInput models.UsersRole
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

// GetRole mengambil role berdasarkan ID atau semua role jika ID tidak disediakan
func GetRole(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		var role models.UsersRole
		if err := models.DB.First(&role, id).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		helper.ResponseJSON(w, http.StatusOK, role)
	} else {
		var roles []models.UsersRole
		if err := models.DB.Find(&roles).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		helper.ResponseJSON(w, http.StatusOK, roles)
	}
}

// UpdateRole mengupdate role berdasarkan ID
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	var roleInput models.UsersRole
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&roleInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	if err := models.DB.Model(&models.UsersRole{}).Where("id = ?", id).Updates(&roleInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Role berhasil diupdate"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// DeleteRole menghapus role berdasarkan ID
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if err := models.DB.Where("id = ?", id).Delete(&models.UsersRole{}).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Role berhasil dihapus"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
