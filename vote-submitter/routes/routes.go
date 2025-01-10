package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/auth"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/handlers"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler handlers.AuthHandler,
	electionHandler handlers.ElectionHandler,
	userHandler handlers.UserHandler,
	voteHandler handlers.VoteHandler) {

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Middleware for all protected routes
	authorized := r.Group("/api")
	{
		authorized.Use(auth.AuthMiddleware()) // Group for authorized routes
	}

	// User routes
	users := authorized.Group("/users")
	{
		users.POST("/", userHandler.CreateUser).Use(auth.RoleMiddleware("admin"))
	}

	// Election routes
	elections := authorized.Group("/elections")
	{
		elections.GET("/", electionHandler.GetElections)
		elections.GET("/:id", electionHandler.GetElection)
		elections.POST("/", electionHandler.CreateElection).Use(auth.RoleMiddleware("admin"))
		elections.PUT("/:id", electionHandler.UpdateElection).Use(auth.RoleMiddleware("admin"))
		elections.POST("/:id/candidates", electionHandler.AddCandidateToElection).Use(auth.RoleMiddleware("admin"))
	}

	// Vote routes
	votes := authorized.Group("/votes")
	{
		votes.POST("/:id", voteHandler.CastVote).Use(auth.RoleMiddleware("voter"))
	}
}
