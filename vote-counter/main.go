package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohsenpakzad/distributed-voting-system/shared/queue"
	"github.com/mohsenpakzad/distributed-voting-system/vote-counter/api"
	"github.com/mohsenpakzad/distributed-voting-system/vote-counter/counter"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	isBootstrap := os.Getenv("BOOTSTRAP") == "true"

	nodeID := os.Getenv("NODE_ID")
	nodeAddress := os.Getenv("NODE_ADDRESS")
	dataDir := os.Getenv("DATA_DIR")
	if nodeID == "" || nodeAddress == "" || dataDir == "" {
		log.Fatal("NODE_ID, NODE_ADDRESS, and DATA_DIR environment variables must be set")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable not set")
	}
	kafkaBrokersArray := strings.Split(kafkaBrokers, ",")

	// Create Raft node
	node, err := counter.NewNode(nodeID, nodeAddress, dataDir, isBootstrap)
	if err != nil {
		log.Fatal(err)
	}

	voteCounter := counter.NewVoteCounter(node)
	voteHandler := queue.NewVoteConsumerHandler(voteCounter)
	validatedVoteConsumer, err := queue.NewVoteConsumer(
		kafkaBrokersArray,
		queue.ValidatedVoteConsumer,
		voteHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create validated vote consumer: %v", err)
	}
	defer validatedVoteConsumer.Close()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming in a goroutine
	go func() {
		if err := validatedVoteConsumer.Start(ctx); err != nil {
			log.Printf("Error starting consumer: %v", err)
			cancel()
		}
	}()

	// Start API server in a goroutine
	go func() {
		router := gin.Default()

		electionResultHandler := api.NewElectionResultHandler(node)

		api.SetupRoutes(router, electionResultHandler)

		port := os.Getenv("API_PORT")
		if port == "" {
			port = "3001"
		}
		router.Run(":" + port)
	}()

	// Wait for shutdown signal
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	<-terminate
}
