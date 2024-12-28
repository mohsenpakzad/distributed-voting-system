package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/auth"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/database"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/utils"
)

func Login(c *gin.Context) {

	db := database.GetDB(c)

	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
