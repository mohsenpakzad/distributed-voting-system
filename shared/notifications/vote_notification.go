package notifications

import (
	"encoding/json"
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

func CreateVoteNotification(db *gorm.DB, vote *models.Vote, title string) {
	var description *string

	// If vote is provided, serialize it into the description field
	if vote != nil {
		// Serialize vote details to JSON
		voteDetails, err := json.Marshal(vote)
		if err == nil {
			// Convert the serialized vote details into a string pointer
			descriptionStr := string(voteDetails)
			description = &descriptionStr
		} else {
			log.Printf("Failed to serialize vote details: %v", err)
		}
	}

	// Create a new notification
	notification := models.Notification{
		UserID:      vote.UserID,
		Title:       title,
		Description: description, // Set description as nil if not provided
		IsRead:      false,
	}

	// Insert the notification into the database
	if err := db.Create(&notification).Error; err != nil {
		log.Printf("Failed to create notification: %v", err)
	}
}
