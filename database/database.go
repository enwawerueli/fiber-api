package database

import (
	"log"

	"github.com/enwawerueli/fiber-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	// Set up database connection
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database.\n%s", err.Error())
	}
	log.Println("Successfuly connected to the database")
	// Run migration
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	DB = db
}
