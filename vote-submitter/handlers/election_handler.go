package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/database"
)

func GetElections(c *gin.Context) {
	db := database.GetDB(c)

	var elections []models.Election
	if err := db.Find(&elections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve elections"})
		return
	}
	c.JSON(http.StatusOK, elections)
}

func GetElection(c *gin.Context) {
	db := database.GetDB(c)

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	var election models.Election
	if err := db.Preload("Candidates").Where("id = ?", id).First(&election).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}
	c.JSON(http.StatusOK, election)
}

func CreateElection(c *gin.Context) {
	db := database.GetDB(c)

	var election models.Election
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.Create(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create election"})
		return
	}

	c.JSON(http.StatusCreated, election)
}

func UpdateElection(c *gin.Context) {
	db := database.GetDB(c)

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
	if err := db.First(&existingElection, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	updatedElection.ID = existingElection.ID // Important: Preserve the ID
	if err := db.Save(&updatedElection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update election"})
		return
	}

	c.JSON(http.StatusOK, updatedElection)
}

func DeleteElection(c *gin.Context) {
	db := database.GetDB(c)

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID format"})
		return
	}

	if err := db.Delete(&models.Election{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete election"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Election deleted successfully"})
}

func AddCandidateToElection(c *gin.Context) {
	db := database.GetDB(c)

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
	candidate.ElectionID = electionID

	if err := db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add candidate"})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}
