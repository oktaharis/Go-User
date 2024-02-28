package models

import (
	"time"
)

// Product adalah model untuk tabel Product
type Product struct {
	ID        uint `gorm:"primaryKey" json:"id"` // Kunci primer
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(50)" json:"name"`
}
