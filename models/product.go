package models



// Product adalah model untuk tabel Product
type Product struct {
	ID        uint `gorm:"primaryKey" json:"id"` // Kunci primer
	Name      string `gorm:"type:varchar(50)" json:"name"`
}
