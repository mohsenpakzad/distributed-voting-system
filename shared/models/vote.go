package models

import (
	"time"
)

type Vote struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // Unique identifier for the vote
	CreatedAt   time.Time `json:"created_at"`                                               // Timestamp for creation
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id" gorm:"index"`      // Identifier of the voter
	ElectionID  string    `json:"election_id" gorm:"index"`  // Identifier of the election
	CandidateID string    `json:"candidate_id" gorm:"index"` // Identifier of the candidate
	Timestamp   time.Time `json:"timestamp"`                 // Time when the vote was cast
}
