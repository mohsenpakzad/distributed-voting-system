package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

type NotificationHandler interface {
	GetAllNotifications(c *gin.Context)
	GetUnreadNotifications(c *gin.Context)
	MarkAsRead(c *gin.Context)
	MarkAllAsRead(c *gin.Context)
}

type notificationHandler struct {
	db *gorm.DB
}

// NewNotificationHandler creates a new NotificationHandler instance.
func NewNotificationHandler(db *gorm.DB) NotificationHandler {
	return &notificationHandler{db}
}

// GetAllNotifications retrieves all notifications for a specific user.
func (h *notificationHandler) GetAllNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var notifications []models.Notification
	if err := h.db.Where("user_id = ?", userID).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadNotifications returns all unread notifications for the authenticated user.
func (h *notificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var notifications []models.Notification
	if err := h.db.Where("user_id = ? AND is_read = false", userID).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Could not fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead marks a single notification as read.
func (h *notificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")

	var notification models.Notification
	// Check if the notification belongs to the user
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		First(&notification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound,
				gin.H{"error": "Notification not found or does not belong to you"})
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "Could not retrieve notification"})
		}
		return
	}

	// Mark the notification as read
	if err := h.db.Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Could not mark notification as read"})
		return
	}

	c.Status(http.StatusOK)
}

// MarkAllAsRead marks all notifications for a specific user as read.
func (h *notificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	if err := h.db.Model(&models.Notification{}).Where("user_id = ?", userID).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Could not mark all notifications as read"})
		return
	}

	c.Status(http.StatusOK)
}
