// models/roles.go

package models

import (
	"time"
)

// Role adalah model untuk tabel Role
type Role struct {
	ID        uint `gorm:"primaryKey" json:"id"` // Kunci primer
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(50)" json:"name"`
}