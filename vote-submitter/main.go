package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/database"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db := database.ConnectDB()

	r := gin.Default()
	r.Use(database.InjectDatabaseMiddleware(db))
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
