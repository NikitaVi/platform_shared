package kafka

import (
	"context"
	"github.com/NikitaVi/platform_shared/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
