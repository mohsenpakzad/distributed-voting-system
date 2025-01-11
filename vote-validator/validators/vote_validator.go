package validators

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
	"gorm.io/gorm"
)

type UnverifiedVoteValidator struct {
	db         *gorm.DB
	validators []VoteValidator
}

func NewUnverifiedVoteValidator(db *gorm.DB) *UnverifiedVoteValidator {
	validators := []VoteValidator{
		NewElectionExistValidator(db),
		NewCandidateExistValidator(db),
		NewDuplicateVoteValidator(db),
	}

	return &UnverifiedVoteValidator{db, validators}
}

// ProcessVote validates and stores the vote in the database
func (p *UnverifiedVoteValidator) ValidateVote(vote models.Vote) {
	// Validate the vote
	for _, validator := range p.validators {
		if err := validator.Validate(vote); err != nil {
			log.Printf("Validation failed: %v", err)
			return
		}
	}

	// Process the valid vote (e.g., save it to the database)
	if err := p.db.Create(&vote).Error; err != nil {
		log.Printf("Failed to save vote: %v", err)
		return
	}

	log.Printf("Vote processed successfully: %v", vote)
}
