package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"github.com/mohsenpakzad/distributed-voting-system/vote-submitter/utils"
	"gorm.io/gorm"
)

type UserHandler interface {
	CreateUser(c *gin.Context);
}

type userHandler struct {
    db *gorm.DB;
}

func NewUserHandler(db *gorm.DB) UserHandler {
    return &userHandler{db}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = hashedPassword

	if err := h.db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
