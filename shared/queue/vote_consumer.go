package queue

import (
	"context"

	"github.com/IBM/sarama"
)

// VoteConsumer handles consuming messages from a Kafka topic
type VoteConsumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler *VoteConsumerHandler
}

type VoteConsumerType struct {
	topic   string
	groupId string
}

var (
	UnverifiedVoteConsumer = VoteConsumerType{unverifiedVoteTopic, unverifiedVotesGroupID}
	ValidatedVoteConsumer  = VoteConsumerType{validatedVoteTopic, validatedVotesGroupID}
)

// NewVoteConsumer creates a new instance of UnverifiedVoteConsumer
func NewVoteConsumer(brokers []string,
	typ VoteConsumerType, handler *VoteConsumerHandler) (*VoteConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokers, typ.groupId, config)
	if err != nil {
		return nil, err
	}

	return &VoteConsumer{
		group:   group,
		topics:  []string{typ.topic},
		handler: handler,
	}, nil
}

// Start begins consuming messages from the Kafka topic
func (c *VoteConsumer) Start(ctx context.Context) error {
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
func (c *VoteConsumer) Close() error {
	return c.group.Close()
}
