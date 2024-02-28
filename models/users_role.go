package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRole struct {
	ID        string `gorm:"type:char(36);primaryKey;" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint `gorm:"int" json:"user_id"`            // Kolom UserID sebagai kunci asing
	User      User `gorm:"foreignKey:UserID" json:"user"` // Kolom UserID sebagai kunci asing
	RoleID    uint `gorm:"int" json:"role_id"`            // Kunci asing untuk Product
	Role      Role `gorm:"foreignKey:RoleID" json:"role"` // Kolom RoleID sebagai kunci as
}

// BeforeCreate will set a UUID rather than numeric ID.
func (up *UsersRole) BeforeCreate(tx *gorm.DB) (err error) {
	uid := uuid.New()
	up.ID = uid.String()
	return nil
}
