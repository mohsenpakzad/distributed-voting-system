package models

import "time"

type Notification struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id" gorm:"index"`
	Title       string    `json:"title"`
	Description *string   `json:"description" gorm:"default:NULL"`
	IsRead      bool      `json:"is_read" gorm:"default:false"`
}
