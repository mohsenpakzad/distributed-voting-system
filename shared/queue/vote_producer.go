package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

// VoteProducer represents a Kafka message producer
type VoteProducer struct {
	producer sarama.SyncProducer
	topic    string
}

type VoteProducerType struct {
	topic string
}

var (
	UnverifiedVoteProducer = VoteProducerType{unverifiedVoteTopic}
	ValidatedVoteProducer  = VoteProducerType{validatedVoteTopic}
)

// NewVoteProducer creates a new Kafka producer
func NewVoteProducer(
	brokers []string,
	typ VoteProducerType) (*VoteProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &VoteProducer{producer, typ.topic}, nil
}

// SendVote sends a vote to the Kafka topic
func (p *VoteProducer) SendVote(vote *models.Vote) error {
	// Serialize the vote to JSON
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return fmt.Errorf("failed to marshal vote: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
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
func (p *VoteProducer) Close() error {
	return p.producer.Close()
}
