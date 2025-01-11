package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

// UnverifiedVoteProducer represents a Kafka message producer
type UnverifiedVoteProducer struct {
	producer sarama.SyncProducer
}

// NewUnverifiedVoteProducer creates a new Kafka producer
func NewUnverifiedVoteProducer(brokers []string) (*UnverifiedVoteProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &UnverifiedVoteProducer{producer: producer}, nil
}

// SendVote sends a vote to the Kafka topic
func (p *UnverifiedVoteProducer) SendVote(vote *models.Vote) error {
	// Serialize the vote to JSON
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return fmt.Errorf("failed to marshal vote: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: unverifiedVoteTopic,
		Key:   sarama.StringEncoder(vote.ID), // Use vote ID as the message key
		Value: sarama.ByteEncoder(voteJSON),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send vote: %w", err)
	}

	log.Printf("Vote sent to partition %d at offset %d\n", partition, offset)
	return nil
}

// Close closes the producer
func (p *UnverifiedVoteProducer) Close() error {
	return p.producer.Close()
}
