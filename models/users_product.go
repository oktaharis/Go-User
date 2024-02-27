package models

import (
	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

// UsersProduct adalah model untuk tabel UsersProduct
type UsersProduct struct {
	gorm.Model
	UserID    string    `gorm:"type:varchar(32);primaryKey;column:id" json:"user_id"` // Kolom UserID sebagai kunci asing
	Name      User      `gorm:"foreignKey:UserID" json:"name"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"` // Kolom ProductID sebagai kunci as
}
