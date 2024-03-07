package rolescontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var roleInput models.Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&roleInput); err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&roleInput).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	data := []models.Role{roleInput}
	response := map[string]interface{}{
		"message": "success",
		"status":  true,
		"data":    data,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func ReadRole(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id dari path parameter
	idParam := mux.Vars(r)["id"]

	// Jika idParam tidak kosong, artinya kita ingin mengambil satu role berdasarkan ID
	if idParam != "" {
		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]interface{}{"message": "Failed, ID tidak valid", "status": false}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data role berdasarkan ID
		var role models.Role
		if err := models.DB.First(&role, id).Error; err != nil {
			response := map[string]interface{}{"message": "Role tidak ditemukan", "status": false}
			helper.ResponseJSON(w, http.StatusNotFound, response)
			return
		}

		// Mengembalikan data role dalam format JSON
		response := map[string]interface{}{"message": "berhasil menampilkan data role", "status": true, "data": role}
		helper.ResponseJSON(w, http.StatusOK, response)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data roles
		var roles []models.Role
		if err := models.DB.Find(&roles).Error; err != nil {
			response := map[string]interface{}{"message": err.Error(), "status": false}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
		// Mengembalikan seluruh data roles dalam format JSON
		if len(roles) > 0 {
			response := map[string]interface{}{"message": "berhasil menampilkan data role", "status": true, "data": roles}
			helper.ResponseJSON(w, http.StatusOK, response)
		}else{
			response := map[string]interface{}{"message": "gagal menampilkan data role", "status": false, "data": roles}
			helper.ResponseJSON(w, http.StatusOK, response)
		}
	}
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
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

	// Mendapatkan data Role yang akan diupdate
	var RoleInput models.Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&RoleInput); err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Cek apakah data Role dengan ID tersebut ada
	var existingRole models.Role
	if err := models.DB.First(&existingRole, id).Error; err != nil {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Update data Role berdasarkan ID
	if err := models.DB.Model(&models.Role{}).Where("id = ?", id).Updates(&RoleInput).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "Role berhasil diupdate", "status": true, "data": existingRole}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
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

	// Cek apakah data Role dengan ID tersebut ada
	var existingRole models.Role
	if err := models.DB.First(&existingRole, id).Error; err != nil {
		response := map[string]interface{}{"message": "Failed, ID tidak ditemukan", "status": false}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Proses penghapusan Role dari database berdasarkan id
	if err := models.DB.Where("id = ?", id).Delete(&models.Role{}).Error; err != nil {
		response := map[string]interface{}{"message": err.Error(), "status": false}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{"message": "Role berhasil dihapus", "status": true}
	helper.ResponseJSON(w, http.StatusOK, response)
}
