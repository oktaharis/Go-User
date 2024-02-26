// models/database.go

package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load() // Load environment variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=dashboard", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Gagal koneksi database")
	} else {
		fmt.Println("Berhasil terhubung dengan PostgreSQL")
	}


	DB = db

        // AutoMigrate akan membuat tabel baru jika belum ada
        db.AutoMigrate(&User{})
        db.AutoMigrate(&Role{})
        db.AutoMigrate(&Product{})
        db.AutoMigrate(&UsersProduct{})
        db.AutoMigrate(&UsersRole{})
}
