package api

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, electionResultHandler ElectionResultHandler) {

	api := router.Group("/api")
	{
		elections := api.Group("/elections-result")
		{
			elections.GET("/:id", electionResultHandler.GetElectionResult)
			elections.GET("/", electionResultHandler.GetElectionsResult)
		}
	}
}
