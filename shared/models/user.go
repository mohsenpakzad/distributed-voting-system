package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // Unique identifier for the user
	CreatedAt time.Time `json:"created_at"`                                               // Timestamp for creation
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Role      string    `json:"role"`
	Password  string
}
