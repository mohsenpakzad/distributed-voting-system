package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/database"
)

func CastVote(c *gin.Context) {
	db := database.GetDB(c)

	electionID := c.Param("id")
	_, err := uuid.Parse(electionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	var input struct {
		CandidateID string `json:"candidate_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	candidateID, err := uuid.Parse(input.CandidateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID format"})
		return
	}

	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	vote := models.Vote{
		UserID:      userID.(string),
		ElectionID:  electionID,
		CandidateID: candidateID.String(),
		Timestamp:   time.Now(),
	}

	if err := db.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast vote"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vote casted successfully"})
}
