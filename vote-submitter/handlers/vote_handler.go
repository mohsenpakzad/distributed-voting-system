package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/shared/queue"
)

type VoteHandler interface {
	CastVote(c *gin.Context)
}

type voteHandler struct {
	unverifiedVoteProducer *queue.VoteProducer
}

func NewVoteHandler(vp *queue.VoteProducer) VoteHandler {
	return &voteHandler{vp}
}

func (h *voteHandler) CastVote(c *gin.Context) {
	var input struct {
		CandidateID string `json:"candidate_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.MustGet("user_id").(string)
	electionID := c.Param("id")

	vote := models.Vote{
		UserID:      userID,
		ElectionID:  electionID,
		CandidateID: input.CandidateID,
		Timestamp:   time.Now(),
	}

	err := h.unverifiedVoteProducer.SendVote(&vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue vote"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vote queued successfully"})
}
