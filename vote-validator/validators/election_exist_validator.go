package validators

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

// ElectionExistValidator checks if the election exists
type ElectionExistValidator struct {
	db *gorm.DB
}

func NewElectionExistValidator(db *gorm.DB) VoteValidator {
	return &ElectionExistValidator{db}
}

func (v *ElectionExistValidator) Validate(vote models.Vote) error {
	var count int64
	err := v.db.Model(&models.Election{}).
		Where("id = ?", vote.ElectionID).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count == 0 {
		return NewValidationError("election does not exist")
	}

	log.Printf("ElectionExistValidator passed")
	return nil
}
