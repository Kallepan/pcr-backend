package database

import (
	"log"

	"gitlab.com/kaka/pcr-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // postgres driver
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {
	Instance, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if dbError != nil {
		log.Fatal("Can not connect to DB: ", dbError)
		panic("Failed to connect to database!")
	}

	log.Println("Connected to database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
}
