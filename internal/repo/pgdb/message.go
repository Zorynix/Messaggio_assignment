package pgdb

import (
	"context"
	"errors"
	"messagio_testsuite/internal/entity"
	repoerrs "messagio_testsuite/internal/repo/repo_errors"
	"messagio_testsuite/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type MessageRepo struct {
	*postgres.Postgres
}

func NewMessageRepo(pg *postgres.Postgres) *MessageRepo {
	return &MessageRepo{pg}
}

func (r *MessageRepo) CreateMessage(ctx context.Context, message entity.Message) (uuid.UUID, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	query := "INSERT INTO messaggio.messages (message) VALUES ($1) RETURNING id"
	var id uuid.UUID
	err = tx.QueryRow(ctx, query, message.Message).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *MessageRepo) GetMessageById(ctx context.Context, id uuid.UUID) (entity.Message, error) {
	query := "SELECT id, message, created_at, processed, processed_at FROM messaggio.messages WHERE id = $1"
	var message entity.Message
	err := r.Pool.QueryRow(ctx, query, id).Scan(&message.ID, &message.Message, &message.CreatedAt, &message.Processed, &message.ProcessedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Message{}, repoerrs.ErrNotFound
		}
		return entity.Message{}, err
	}
	return message, nil
}

func (r *MessageRepo) GetMessages(ctx context.Context) ([]entity.Message, error) {
	query := "SELECT id, message, created_at, processed, processed_at FROM messaggio.messages"
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var message entity.Message
		if err := rows.Scan(&message.ID, &message.Message, &message.CreatedAt, &message.Processed, &message.ProcessedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepo) MarkMessageAsProcessed(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE messaggio.messages SET processed = true, processed_at = CURRENT_TIMESTAMP WHERE id = $1"
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *MessageRepo) GetProcessedMessagesStats(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM messaggio.messages WHERE processed = true"
	var count int
	err := r.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *MessageRepo) GetMessageByContent(ctx context.Context, content string) (entity.Message, error) {
	query := "SELECT id, message, created_at, processed, processed_at FROM messaggio.messages WHERE message = $1"
	var message entity.Message
	err := r.Pool.QueryRow(ctx, query, content).Scan(&message.ID, &message.Message, &message.CreatedAt, &message.Processed, &message.ProcessedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Message{}, repoerrs.ErrNotFound
		}
		return entity.Message{}, err
	}
	return message, nil
}
