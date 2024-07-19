package repo

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo/pgdb"
	"messagio_testsuite/pkg/postgres"

	"github.com/google/uuid"
)

type Message interface {
	CreateMessage(ctx context.Context, message entity.Message) (uuid.UUID, error)
	GetMessageById(ctx context.Context, id uuid.UUID) (entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
	MarkMessageAsProcessed(ctx context.Context, id uuid.UUID) error
	GetProcessedMessagesStats(ctx context.Context) (int, error)
}

type Repositories struct {
	Message
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Message: pgdb.NewMessageRepo(pg),
	}
}
