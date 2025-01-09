package database

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dbUrl string) *gorm.DB {
	var err error
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Election{}, &models.Candidate{}, &models.Vote{})
	if err != nil {
		log.Fatal("failed to automigrate:", err)
	}

	return db
}
