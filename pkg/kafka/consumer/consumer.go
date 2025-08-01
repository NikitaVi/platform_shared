package consumer

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
	"strings"
)

type consumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
) *consumer {
	return &consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

func (c *consumer) Consume(ctx context.Context, topicName string, handler Handler) error {
	c.consumerGroupHandler.msgHandler = handler
	return c.consume(ctx, topicName)
}

func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *consumer) consume(ctx context.Context, topicName string) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(topicName, ","), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("rebalancing...\n")
	}
}
