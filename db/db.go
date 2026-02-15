package db

import (
	"log"
	"os"

	"APPDROP/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal(" DATABASE_URL is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	DB = database
	log.Println("Connected to PostgreSQL")

	_ = DB.Exec(`ALTER TABLE pages DROP CONSTRAINT IF EXISTS pages_route_key`).Error
	_ = DB.Exec(`ALTER TABLE pages DROP CONSTRAINT IF EXISTS uni_pages_route`).Error
	_ = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_pages_brand_route ON pages(brand_id, route)`).Error

	if err := DB.AutoMigrate(&models.Page{}, &models.Widget{}); err != nil {
		log.Println("Failed to migrate Pages/Widgets:", err)
	}

	if err := DB.Exec(`ALTER TABLE brands ADD COLUMN IF NOT EXISTS email text;`).Error; err != nil {
		log.Println("Failed to add email column:", err)
	}

	if err := DB.Exec(`ALTER TABLE brands ADD COLUMN IF NOT EXISTS password_hash text;`).Error; err != nil {
		log.Println("Failed to add password_hash column:", err)
	}
}
