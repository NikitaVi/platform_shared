package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type GroupHandler struct {
	msgHandler Handler
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

func (g *GroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (g *GroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (g *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("message channel closed")
				return nil
			}

			log.Printf("message claimed value: %s, topic: %s, partition: %d\n", message.Value, message.Topic, message.Partition)
			err := g.msgHandler(session.Context(), message)
			if err != nil {
				log.Printf("failed to handle message: %v", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Println("session closed")
			return nil
		}
	}
}
