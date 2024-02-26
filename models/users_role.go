package models

import (
	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

// UsersRole adalah model untuk tabel UsersRole
type UsersRole struct {
	gorm.Model
	UserID  string    `gorm:"type:varchar(32);primaryKey;column:id" json:"user_id"`
	RoleID  uuid.UUID `gorm:"type:uuid" json:"role_id"`
}
