package config

import (
	"log"
	"navi-wings/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	dsn := "host=localhost user=postgres password=password dbname=navi_wings_dev port=5432 sslmode=disable"

	// Log the connection string for debugging
	log.Println("Connecting to database with DSN:", dsn)

	// Open the database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to the database!")

	// Run AutoMigrate
	if err := DB.AutoMigrate(&models.Message{}); err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}

	log.Println("Auto migration complete!")
}
