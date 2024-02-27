package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersProduct struct {
	ID        string  `gorm:"type:char(36);primaryKey;" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int     `gorm:"int" json:"user_id"` // Kolom UserID sebagai kunci asing
	User      User    `gorm:"foreignKey:UserID" json:"user"` // Kolom UserID sebagai kunci asing
	ProductID uint    `gorm:"int" json:"product_id"`               // Kunci asing untuk Product
	Product   Product `gorm:"foreignKey:ProductID" json:"product"` // Kolom ProductID sebagai kunci as
}



// BeforeCreate will set a UUID rather than numeric ID.
func (up *UsersProduct) BeforeCreate(tx *gorm.DB) (err error) {
	uid := uuid.New()
	up.ID = uid.String()
	return nil
}
