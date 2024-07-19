package service

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo"
	"messagio_testsuite/pkg/kafka"

	"github.com/google/uuid"
)

type Message interface {
	CreateMessage(ctx context.Context, content string) (uuid.UUID, error)
	GetMessageById(ctx context.Context, messageId uuid.UUID) (entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
	MarkMessageAsProcessed(ctx context.Context, messageId uuid.UUID) error
	GetProcessedMessagesStats(ctx context.Context) (int, error)
}

type Services struct {
	Message Message
}

type ServicesDependencies struct {
	Repos         *repo.Repositories
	KafkaProducer *kafka.KafkaProducer
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Message: NewMessageService(deps.Repos.Message, deps.KafkaProducer),
	}
}
