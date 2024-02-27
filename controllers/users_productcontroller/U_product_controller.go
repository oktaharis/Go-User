package usersproductcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
	"gorm.io/gorm"
)

func CreateUserProduct(w http.ResponseWriter, r *http.Request) {
	var UproductInput models.UsersProduct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&UproductInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// INSERT KE DATABASE
	if err := models.DB.Create(&UproductInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	response := map[string]string{"message": "success"}
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
			response := map[string]string{"message": "ID tidak valid"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		// Mendapatkan data UserProduct berdasarkan ID
		var UserProduct models.UsersProduct
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Product").Select("id, user_id, product_id").First(&UserProduct, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response := map[string]string{"message": "User tidak ditemukan"}
				helper.ResponseJSON(w, http.StatusNotFound, response)
				return
			}
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan data User dalam format JSON tanpa CreatedAt
		userWithoutCreatedAt := struct {
			ID        string      `json:"id"`
			UserID    int         `json:"user_id"`
			User      models.User `json:"user"`
			ProductID uint        `json:"product_id"`
			Product   models.Product `json:"product"`
		}{
			ID:        UserProduct.ID,
			UserID:    UserProduct.UserID,
			User:      UserProduct.User,
			ProductID: UserProduct.ProductID,
			Product:   UserProduct.Product,
		}

		helper.ResponseJSON(w, http.StatusOK, userWithoutCreatedAt)
	} else {
		// Jika idParam kosong, artinya kita ingin mengambil seluruh data Users
		var users []models.UsersProduct
		if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).Preload("Product").Select("id, user_id, product_id").Find(&users).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}

		// Mengembalikan seluruh data Users dalam format JSON tanpa CreatedAt
		var usersWithoutCreatedAt []struct {
			ID        string      `json:"id"`
			UserID    int         `json:"user_id"`
			User      models.User `json:"user"`
			ProductID uint        `json:"product_id"`
			Product   models.Product `json:"product"`
		}
		for _, user := range users {
			userWithoutCreatedAt := struct {
				ID        string      `json:"id"`
				UserID    int         `json:"user_id"`
				User      models.User `json:"user"`
				ProductID uint        `json:"product_id"`
				Product   models.Product `json:"product"`
			}{
				ID:        user.ID,
				UserID:    user.UserID,
				User:      user.User,
				ProductID: user.ProductID,
				Product:   user.Product,
			}
			usersWithoutCreatedAt = append(usersWithoutCreatedAt, userWithoutCreatedAt)
		}

		helper.ResponseJSON(w, http.StatusOK, usersWithoutCreatedAt)
	}
}




