//models users_product
package models

import (
	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

// UsersProduct adalah model untuk tabel UsersProduct
type UsersProduct struct {
	gorm.Model
	UserID    string    `gorm:"type:varchar(32);primaryKey;column:id" json:"user_id"` // Kolom UserID sebagai kunci asing
	Name      string    `gorm:"type:varchar(255)" json:"name"`                         // Kolom Name
	Code      string    `gorm:"type:varchar(50)" json:"code"`                          // Kolom Code
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`                          // Kolom ProductID sebagai kunci as
}

