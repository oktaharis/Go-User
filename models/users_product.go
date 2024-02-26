package models

import (
	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

// UsersProduct adalah model untuk tabel UsersProduct
type UsersProduct struct {
	gorm.Model
	UserID    string    `gorm:"type:varchar(32);primaryKey;column:id" json:"user_id"` // Kolom UserID sebagai kunci asing
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`                          // Kolom ProductID sebagai kunci as
}
