package pgdb

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/pkg/postgres"
)

type MessageRepo struct {
	*postgres.Postgres
}

func NewMessageRepo(pg *postgres.Postgres) *MessageRepo {
	return &MessageRepo{pg}
}

func (r *MessageRepo) CreateMessage(ctx context.Context) (int, error) {
	return 0, nil
}

func (r *MessageRepo) GetMessageById(ctx context.Context, id int) (entity.Message, error) {

	var message entity.Message

	return message, nil
}

func (r *MessageRepo) GetMessages(ctx context.Context) ([]entity.Message, error) {

	var messages []entity.Message

	return messages, nil
}
