package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

type ElectionHandler interface {
	GetElections(c *gin.Context);
	GetElection(c *gin.Context);
	CreateElection(c *gin.Context);
	UpdateElection(c *gin.Context);
	AddCandidateToElection(c *gin.Context);
}

type electionHandler struct {
    db *gorm.DB;
}

func NewElectionHandler(db *gorm.DB) ElectionHandler {
    return &electionHandler{db}
}

func (h *electionHandler)GetElections(c *gin.Context) {
	var elections []models.Election
	if err := h.db.Preload("Candidates").Find(&elections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve elections"})
		return
	}
	c.JSON(http.StatusOK, elections)
}

func (h *electionHandler) GetElection(c *gin.Context) {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	var election models.Election
	if err := h.db.Preload("Candidates").
		Where("id = ?", id).
		First(&election).Error;
		err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}
	c.JSON(http.StatusOK, election)
}

func (h *electionHandler) CreateElection(c *gin.Context) {
	var election models.Election
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.db.Create(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create election"})
		return
	}

	c.JSON(http.StatusCreated, election)
}

func (h *electionHandler) UpdateElection(c *gin.Context) {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	var updatedElection models.Election
	if err := c.ShouldBindJSON(&updatedElection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var existingElection models.Election
	if err := h.db.First(&existingElection, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	updatedElection.ID = existingElection.ID // Important: Preserve the ID
	if err := h.db.Save(&updatedElection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update election"})
		return
	}

	c.JSON(http.StatusOK, updatedElection)
}

func (h *electionHandler) AddCandidateToElection(c *gin.Context) {
	electionID := c.Param("id")
	_, err := uuid.Parse(electionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	var candidate models.Candidate
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var election models.Election
	if err := h.db.First(&election, "id = ?", electionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	candidate.ElectionID = electionID

	if err := h.db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add candidate"})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}
