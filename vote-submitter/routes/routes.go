package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/auth"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/handlers"
)

func SetupRoutes(r *gin.Engine) {

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", handlers.Login)
	}

	// Middleware for all protected routes
	authorized := r.Group("/api")
	{
		authorized.Use(auth.AuthMiddleware()) // Group for authorized routes
	}

	// User routes
	users := authorized.Group("/users")
	{
		users.POST("/", handlers.CreateUser).Use(auth.RoleMiddleware("admin"))
	}

	// Election routes
	elections := authorized.Group("/elections")
	{
		elections.GET("/", handlers.GetElections)
		elections.GET("/:id", handlers.GetElection)
		elections.POST("/", handlers.CreateElection).Use(auth.RoleMiddleware("admin"))
		elections.PUT("/:id", handlers.UpdateElection).Use(auth.RoleMiddleware("admin"))
		elections.DELETE("/:id", handlers.DeleteElection).Use(auth.RoleMiddleware("admin"))
		elections.POST("/:id/candidates", handlers.AddCandidateToElection).Use(auth.RoleMiddleware("admin"))
	}

	// Vote routes
	votes := authorized.Group("/votes")
	{
		votes.POST("/:id", handlers.CastVote).Use(auth.RoleMiddleware("voter"))
	}
}
