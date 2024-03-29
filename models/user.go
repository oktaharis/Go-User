// models/models.go

package models

import (
	"time"
)

type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UID             string 	  `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string    `gorm:"varchar(255)" json:"name"`
	Email           string    `gorm:"varchar(255)" json:"email"`
	EmailVerifiedAt time.Time `gorm:"timestamp" json:"email_verified_at"`
	Password        string    `gorm:"varchar(255)" json:"password"`
	Phone           string    `gorm:"varchar(255)" json:"phone"`
	LastLogin       time.Time `gorm:"timestamp" json:"last_login"`
	Status          string    `gorm:"varchar(100)" json:"status"`
	RememberToken   string    `gorm:"varchar(100)" json:"remember_token"`
	RoleID          uint      `gorm:"int" json:"role_id"` // Kunci asing untuk Role
	Role            Role      `gorm:"foreignKey:RoleID" json:"role"`
	ProductID       uint      `gorm:"int" json:"product_id"` // Kunci asing untuk Product
	Product         Product   `gorm:"foreignKey:ProductID" json:"product"`
	OTP             string    `gorm:"varchar(6)" json:"otp"`
	ExpiredAt       time.Time `gorm:"timestamp" json:"expired_at"`
	ResetPassword   string    `gorm:"varchar(10)" json:"new_password"`
}
