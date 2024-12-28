package models

import "time"

type Candidate struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ElectionID  string    `json:"election_id" gorm:"index"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
