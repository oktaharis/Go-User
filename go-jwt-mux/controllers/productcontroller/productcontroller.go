package productcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
)

func CreateProduct(w http.ResponseWriter, r *http.Request){
	var productInput models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&productInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

		// INSERT KE DATABASE
		if err := models.DB.Create(&productInput).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
	
		response := map[string]string{"message": "success"}
		helper.ResponseJSON(w, http.StatusOK, response)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
		// Mendapatkan parameter id dari query parameter
		idParam := r.URL.Query().Get("id")

		// Jika idParam tidak kosong, artinya kita ingin mengambil satu Product berdasarkan ID
		if idParam != "" {
			// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
			id, err := strconv.Atoi(idParam)
			if err != nil {
				response := map[string]string{"message": "ID tidak valid"}
				helper.ResponseJSON(w, http.StatusBadRequest, response)
				return
			}
	
			// Mendapatkan data Product berdasarkan ID
			var product models.Product
			if err := models.DB.First(&product, id).Error; err != nil {
				response := map[string]string{"message": err.Error()}
				helper.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}
	
			// Mengembalikan data Product dalam format JSON
			helper.ResponseJSON(w, http.StatusOK, product)
		} else {
			// Jika idParam kosong, artinya kita ingin mengambil seluruh data Product
			var product []models.Product
			if err := models.DB.Find(&product).Error; err != nil {
				response := map[string]string{"message": err.Error()}
				helper.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}
	
			// Mengembalikan seluruh data Product dalam format JSON
			helper.ResponseJSON(w, http.StatusOK, product)
		}
}
func UpdateProduct(w http.ResponseWriter, r *http.Request){
	// Mendapatkan parameter id dari query parameter
	idParam := r.URL.Query().Get("id")

	// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"message": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Mendapatkan data Product yang akan diupdate
	var productInput models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&productInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	// Update data Product berdasarkan ID
	if err := models.DB.Model(&models.Product{}).Where("id = ?", id).Updates(&productInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Product berhasil diupdate"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func DeleteProduct(w http.ResponseWriter, r *http.Request){
		// Mendapatkan parameter id dari query parameter
		idParam := r.URL.Query().Get("id")

		// Konversi idParam menjadi tipe data yang sesuai (misalnya, integer)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
	
		// Proses penghapusan Product dari database berdasarkan id
		if err := models.DB.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	
		response := map[string]string{"message": "Product berhasil dihapus"}
		helper.ResponseJSON(w, http.StatusOK, response)
}