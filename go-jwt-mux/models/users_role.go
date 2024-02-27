package models

import (
    "gorm.io/gorm"
    "github.com/google/uuid"
)

// UsersRole adalah model untuk tabel UsersRole
type UsersRole struct {
    gorm.Model
    User_ID   string    `gorm:"type:varchar(32);primaryKey" json:"user_id"`
    User_Name string    `gorm:"type:varchar(255)" json:"user_name"`
    User_Code string    `gorm:"type:varchar(50)" json:"user_code"`
    Role_ID   uuid.UUID `gorm:"type:uuid" json:"role_id"`
    Role_Name string    `gorm:"type:varchar(200)" json:"role_name"`
}
