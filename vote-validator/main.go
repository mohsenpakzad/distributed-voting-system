package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mohsenpakzad/distributed-voting-system/shared/database"
	"github.com/mohsenpakzad/distributed-voting-system/shared/queue"
	"github.com/mohsenpakzad/distributed-voting-system/vote-validator/validators"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	db := database.ConnectDB(dbUrl)
	defer database.CloseDB(db)

	// Kafka configuration
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable not set")
	}

	// Set up the vote processor implementation
	UnverifiedVoteValidator := validators.NewUnverifiedVoteValidator(db)

	// Set up the consumer handler with the vote processor injected
	voteHandler := queue.NewVoteConsumerHandler(UnverifiedVoteValidator)

	unverifiedVoteConsumer, err := queue.NewVoteConsumer(
		strings.Split(kafkaBrokers, ","),
		queue.UnverifiedVoteConsumer,
		voteHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create unverified vote consumer: %v", err)
	}
	defer unverifiedVoteConsumer.Close()

	// Context and graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigchan
		cancel()
	}()

	log.Println("Starting unverified vote consumer...")
	if err := unverifiedVoteConsumer.Start(ctx); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	log.Println("Consumer shut down gracefully.")
}
