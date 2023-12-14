package database

import (
	"log"
	"os"

	"github.com/abdulmanafc2001/gigahive/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("successfully connected to database: %v \n", DB.Name())
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Freelancer{})
	DB.AutoMigrate(&models.Bid{})
}
