package service

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo"
	repoerrs "messagio_testsuite/internal/repo/repo_errors"
	serviceerrs "messagio_testsuite/internal/service/service_errors"
)

type MessageService struct {
	messageRepo repo.Message
}

func NewMessageService(messageRepo repo.Message) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context) (int, error) {
	id, err := s.messageRepo.CreateMessage(ctx)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return 0, serviceerrs.ErrMessageAlreadyExists
		}
		return 0, serviceerrs.ErrCannotCreateMessage
	}

	return id, nil
}
func (s *MessageService) GetMessageById(ctx context.Context, messageId int) (entity.Message, error) {
	return s.messageRepo.GetMessageById(ctx, messageId)
}
func (s *MessageService) GetMessages(ctx context.Context) ([]entity.Message, error) {
	return s.messageRepo.GetMessages(ctx)
}
