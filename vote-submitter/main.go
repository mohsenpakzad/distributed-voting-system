package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohsenpakzad/distributed-voting-system/shared/database"
	"github.com/mohsenpakzad/distributed-voting-system/shared/queue"
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
	defer database.CloseDB(db)

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable not set")
	}

	unverifiedVoteProducer, err := queue.NewVoteProducer(
		strings.Split(kafkaBrokers, ","),
		queue.UnverifiedVoteProducer,
	)
	if err != nil {
		log.Fatalf("Failed to create unverified vote producer: %v", err)
	}
	defer unverifiedVoteProducer.Close()

	authHandler := handlers.NewAuthHandler(db)
	electionHandler := handlers.NewElectionHandler(db)
	userHandler := handlers.NewUserHandler(db)
	voteHandler := handlers.NewVoteHandler(unverifiedVoteProducer)
	notificationHandler := handlers.NewNotificationHandler(db)

	r := gin.Default()
	routes.SetupRoutes(r,
		authHandler,
		electionHandler,
		userHandler,
		voteHandler,
		notificationHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
