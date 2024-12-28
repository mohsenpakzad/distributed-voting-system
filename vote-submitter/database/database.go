package database

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var err error
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Election{}, &models.Candidate{}, &models.Vote{})
	if err != nil {
		log.Fatal("failed to automigrate:", err)
	}

	return db
}

func GetDB(c *gin.Context) *gorm.DB {
	db, exists := c.Get("db")
	if !exists {
		panic("Database connection not found in context") // Or handle appropriately
	}
	return db.(*gorm.DB)
}
