package repo

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo/pgdb"
	"messagio_testsuite/pkg/postgres"
)

type Message interface {
	CreateMessage(ctx context.Context) (int, error)
	GetMessageById(ctx context.Context, id int) (entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
}

type Repositories struct {
	Message
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Message: pgdb.NewMessageRepo(pg),
	}
}
