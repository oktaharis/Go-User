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
	UserID    int  `gorm:"int" json:"user_id"` // Kolom UserID sebagai kunci asing
	Name      User `gorm:"foreignKey:UserID" json:"name"`
	RoleID    uint `gorm:"int" json:"role_id"`            // Kunci asing untuk Role
	Role      Role `gorm:"foreignKey:RoleID" json:"Role"` // Kolom RoleID sebagai kunci as
}

// BeforeCreate untuk menghasilkan UUID baru sebelum membuat entitas baru.
func (up *UsersRole) BeforeCreate(tx *gorm.DB) (err error) {
	uid := uuid.New()
	up.ID = uid.String()
	return nil
}

