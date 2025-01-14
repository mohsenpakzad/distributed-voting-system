package counter

import (
	"log"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

type VoteCounter struct {
	node *Node
}

func NewVoteCounter(node *Node) *VoteCounter {
	return &VoteCounter{node}
}

func (p *VoteCounter) ProcessVote(vote models.Vote) {
	// Add the vote to the Raft cluster
	if err := p.node.AddVote(&vote); err != nil {
		// Handle error (log it, maybe retry, etc.)
		log.Printf("Failed to add vote to Raft cluster: %v", err)
	}
}
