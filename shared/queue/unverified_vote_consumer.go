package queue

import (
	"context"

	"github.com/IBM/sarama"
)

// UnverifiedVoteConsumer handles consuming messages from a Kafka topic
type UnverifiedVoteConsumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler *VoteConsumerHandler
}

// NewUnverifiedVoteConsumer creates a new instance of UnverifiedVoteConsumer
func NewUnverifiedVoteConsumer(brokers []string, handler *VoteConsumerHandler) (*UnverifiedVoteConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokers, unverifiedVotesGroupID, config)
	if err != nil {
		return nil, err
	}

	return &UnverifiedVoteConsumer{
		group:   group,
		topics:  []string{unverifiedVoteTopic},
		handler: handler,
	}, nil
}

// Start begins consuming messages from the Kafka topic
func (c *UnverifiedVoteConsumer) Start(ctx context.Context) error {
	for {
		err := c.group.Consume(ctx, c.topics, c.handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

// Close shuts down the consumer group
func (c *UnverifiedVoteConsumer) Close() error {
	return c.group.Close()
}
