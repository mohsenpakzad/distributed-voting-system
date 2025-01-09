package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/shared/queue"
	"gorm.io/gorm"
)

type VoteHandler interface{
	CastVote(c *gin.Context);
}

type voteHandler struct {
    db *gorm.DB;
}

func NewVoteHandler(db *gorm.DB) VoteHandler {
    return &voteHandler{db}
}

func (h *voteHandler) CastVote(c *gin.Context) {
	producer := queue.GetKafkaProducer()

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

	// Serialize the vote to JSON
	message, err := json.Marshal(vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize vote"})
		return
	}

	// Publish the vote using Sarama
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "votes",
		Value: sarama.StringEncoder(message),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue vote"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vote queued successfully"})
}
