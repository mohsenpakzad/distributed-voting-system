package models

import (
	"time"
)

type Election struct {
	ID          string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Name        string      `json:"name"`
	Candidates  []Candidate `json:"candidates" gorm:"foreignKey:ElectionID"`
	Description string      `json:"description"`
}
