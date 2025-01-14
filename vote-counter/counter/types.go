package counter

import "sync"

// Command represents operations that can be performed on the FSM
type Command struct {
	Op          string `json:"op"`
	ElectionID  string `json:"election_id"`
	CandidateID string `json:"candidate_id"`
}

// VoteState maintains the state of votes
type VoteState struct {
	// ElectionID -> CandidateID -> Count
	Elections map[string]map[string]int
	mu        sync.RWMutex
}
