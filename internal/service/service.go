package service

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo"
)

type Message interface {
	CreateMessage(ctx context.Context) (int, error)
	GetMessageById(ctx context.Context, messageId int) (entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
}

type Services struct {
	Message Message
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func newServices(deps ServicesDependencies) *Services {
	return &Services{
		Message: NewMessageService(deps.Repos.Message),
	}
}
