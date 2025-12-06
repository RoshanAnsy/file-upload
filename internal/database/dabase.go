// package database

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func ConnectDatabase() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	dsn := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
// 		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
// 	)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	DB = db
// 	fmt.Println("✅ Connected to PostgreSQL database successfully!")
// }



package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	_ = godotenv.Load() // Safe load

	// Try DATABASE_URL first
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set in environment variables")
	}

	// Connect directly using DATABASE_URL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL successfully!")
}
