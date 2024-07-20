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
	"github.com/sirupsen/logrus"
)

type MessageService struct {
	messageRepo   repo.Message
	kafkaProducer *kafka.KafkaProducer
	kafkaConsumer *kafka.KafkaConsumer
}

func NewMessageService(messageRepo repo.Message, kafkaProducer *kafka.KafkaProducer, kafkaConsumer *kafka.KafkaConsumer) *MessageService {
	s := &MessageService{
		messageRepo:   messageRepo,
		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,
	}

	go s.listenForProcessedMessages()

	return s
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

func (s *MessageService) GetMessageByContent(ctx context.Context, content string) (entity.Message, error) {
	message, err := s.messageRepo.GetMessageByContent(ctx, content)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Message{}, serviceerrs.ErrMessageNotFound
		}
		return entity.Message{}, err
	}
	return message, nil
}

func (s *MessageService) MarkMessageAsProcessed(ctx context.Context, messageId uuid.UUID) error {
	return s.messageRepo.MarkMessageAsProcessed(ctx, messageId)
}

func (s *MessageService) GetProcessedMessagesStats(ctx context.Context) (int, error) {
	return s.messageRepo.GetProcessedMessagesStats(ctx)
}

func (s *MessageService) listenForProcessedMessages() {
	ctx := context.Background()
	s.kafkaConsumer.Consume(ctx, s.handleMessage)
}

func (s *MessageService) handleMessage(ctx context.Context, messageContent string) {
	message, err := s.GetMessageByContent(ctx, messageContent)
	if err != nil {
		logrus.Errorf("Failed to get message by content: %v", err)
		return
	}

	err = s.MarkMessageAsProcessed(ctx, message.ID)
	if err != nil {
		logrus.Errorf("Failed to mark message %s as processed: %v", message.ID, err)
	} else {
		logrus.Infof("Message %s marked as processed", message.ID)
	}
}
