package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohsenpakzad/distributed-voting-system/shared/database"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/handlers"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	db := database.ConnectDB(dbUrl)

	authHandler := handlers.NewAuthHandler(db)
	electionHandler := handlers.NewElectionHandler(db)
	userHandler := handlers.NewUserHandler(db)
	voteHandler := handlers.NewVoteHandler(db)

	r := gin.Default()
	routes.SetupRoutes(r, authHandler, electionHandler, userHandler, voteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
