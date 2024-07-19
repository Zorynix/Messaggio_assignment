package service

import (
	"context"
	"errors"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo"
	repoerrs "messagio_testsuite/internal/repo/repo_errors"
	serviceerrs "messagio_testsuite/internal/service/service_errors"
	"messagio_testsuite/pkg/kafka"

	"github.com/google/uuid"
)

type MessageService struct {
	messageRepo   repo.Message
	kafkaProducer *kafka.KafkaProducer
}

func NewMessageService(messageRepo repo.Message, kafkaProducer *kafka.KafkaProducer) *MessageService {
	return &MessageService{
		messageRepo:   messageRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, content string) (uuid.UUID, error) {
	message := entity.Message{
		Message: content,
	}

	id, err := s.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return uuid.Nil, serviceerrs.ErrCannotCreateMessage
	}

	err = s.kafkaProducer.Produce(content)
	if err != nil {
		return uuid.Nil, serviceerrs.ErrCannotProduceMessage
	}

	return id, nil
}

func (s *MessageService) GetMessageById(ctx context.Context, messageId uuid.UUID) (entity.Message, error) {
	message, err := s.messageRepo.GetMessageById(ctx, messageId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Message{}, serviceerrs.ErrMessageNotFound
		}
		return entity.Message{}, err
	}
	return message, nil
}

func (s *MessageService) GetMessages(ctx context.Context) ([]entity.Message, error) {
	return s.messageRepo.GetMessages(ctx)
}

func (s *MessageService) MarkMessageAsProcessed(ctx context.Context, messageId uuid.UUID) error {
	return s.messageRepo.MarkMessageAsProcessed(ctx, messageId)
}

func (s *MessageService) GetProcessedMessagesStats(ctx context.Context) (int, error) {
	return s.messageRepo.GetProcessedMessagesStats(ctx)
}
