package validators

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

// DuplicateVoteValidator checks if a user has already voted
type DuplicateVoteValidator struct {
	db *gorm.DB
}

func NewDuplicateVoteValidator(db *gorm.DB) VoteValidator {
	return &DuplicateVoteValidator{db}
}

func (v *DuplicateVoteValidator) Validate(vote models.Vote) error {
	var count int64
	err := v.db.Model(&models.Vote{}).
		Where("user_id = ? AND election_id = ?", vote.UserID, vote.ElectionID).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count > 0 {
		return NewValidationError("user has already voted in this election")
	}

	log.Printf("DuplicateVoteValidator passed")
	return nil
}
