package queue

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

// VoteConsumerHandler handles the messages consumed by the Kafka consumer group
type VoteConsumerHandler struct {
	voteProcessor VoteProcessor
	ready         chan bool
}

// VoteProcessor defines the interface for processing votes
type VoteProcessor interface {
	ProcessVote(vote models.Vote)
}

// NewVoteConsumerHandler creates a new UnverifiedVoteConsumerHandler with the provided VoteProcessor
func NewVoteConsumerHandler(voteProcessor VoteProcessor) *VoteConsumerHandler {
	return &VoteConsumerHandler{
		voteProcessor: voteProcessor,
		ready:         make(chan bool),
	}
}

// Setup is called before starting to consume messages
func (h *VoteConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

// Cleanup is called after consuming messages
func (h *VoteConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.ready = make(chan bool)
	return nil
}

// ConsumeClaim processes each message from the Kafka topic
func (h *VoteConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var vote models.Vote
		err := json.Unmarshal(message.Value, &vote)
		if err != nil {
			log.Printf("Error unmarshaling vote: %v", err)
			continue
		}

		log.Printf("Received vote: %+v\n", vote)

		h.voteProcessor.ProcessVote(vote)

		// Mark the message as processed
		session.MarkMessage(message, "")
	}
	return nil
}
