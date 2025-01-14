package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/vote-counter/counter"
)

type ElectionResultHandler interface {
	GetElectionResult(c *gin.Context)
	GetElectionsResult(c *gin.Context)
}

type electionResultHandler struct {
	node *counter.Node
}

func NewElectionResultHandler(node *counter.Node) ElectionResultHandler {
	return &electionResultHandler{node}
}

func (h *electionResultHandler) GetElectionResult(c *gin.Context) {
	electionID := c.Param("id")

	result := h.node.GetResults(electionID)
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Election not found",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *electionResultHandler) GetElectionsResult(c *gin.Context) {
	results := h.node.GetAllResults()
	c.JSON(http.StatusOK, results)
}
