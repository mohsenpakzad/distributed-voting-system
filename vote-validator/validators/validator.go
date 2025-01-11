package validators

import (
	"fmt"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

// VoteValidator defines the interface for vote validation
type VoteValidator interface {
	Validate(vote models.Vote) error
}

// VoteValidationError represents a custom error for vote validation failures
type VoteValidationError struct {
	Reason string
}

func (e *VoteValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Reason)
}

// NewValidationError creates a new ValidationError
func NewValidationError(reason string) error {
	return &VoteValidationError{Reason: reason}
}
