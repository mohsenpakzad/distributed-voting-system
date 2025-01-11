package validators

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

// CandidateExistValidator checks if the candidate exists in the election
type CandidateExistValidator struct {
	db *gorm.DB
}

func NewCandidateExistValidator(db *gorm.DB) VoteValidator {
	return &CandidateExistValidator{db}
}

func (v *CandidateExistValidator) Validate(vote models.Vote) error {
	var count int64
	err := v.db.Model(&models.Candidate{}).
		Where("id = ? AND election_id = ?", vote.CandidateID, vote.ElectionID).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count == 0 {
		return NewValidationError("candidate does not exist in the election")
	}

	log.Printf("CandidateExistValidator passed")
	return nil
}
