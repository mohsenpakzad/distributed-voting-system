package models

type ElectionResult struct {
	ElectionID string           `json:"election_id"`
	Candidates []CandidateCount `json:"candidates"`
}
