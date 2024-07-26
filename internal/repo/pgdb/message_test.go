package pgdb_test

import (
	"context"
	"messagio_testsuite/internal/entity"
	"messagio_testsuite/internal/repo/pgdb"
	"messagio_testsuite/pkg/postgres"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *postgres.Postgres

const createSchemaSQL = `
CREATE SCHEMA IF NOT EXISTS messaggio;
CREATE TABLE messaggio.messages (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP
);
`

func setupPostgres(t *testing.T) func() {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "messaggio",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port("5432/tcp")),
		),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, nat.Port("5432"))
	require.NoError(t, err)

	dsn := "postgres://postgres:postgres@" + host + ":" + port.Port() + "/messaggio?sslmode=disable"
	pg, err := postgres.New(dsn)
	require.NoError(t, err)

	_, err = pg.Pool.Exec(ctx, createSchemaSQL)
	require.NoError(t, err)

	testDB = pg

	return func() {
		pg.Close()
		postgresC.Terminate(ctx)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestMessageRepo_CreateMessage(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	message := entity.Message{
		Message: "test message",
	}

	id, err := repo.CreateMessage(ctx, message)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestMessageRepo_GetMessageById(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	message := entity.Message{
		Message: "test message",
	}

	id, err := repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	fetchedMessage, err := repo.GetMessageById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, message.Message, fetchedMessage.Message)
	assert.Equal(t, id, fetchedMessage.ID)
}

func TestMessageRepo_GetMessages(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	_, err := repo.CreateMessage(ctx, entity.Message{Message: "test message 1"})
	require.NoError(t, err)

	_, err = repo.CreateMessage(ctx, entity.Message{Message: "test message 2"})
	require.NoError(t, err)

	messages, err := repo.GetMessages(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(messages), 2)
}

func TestMessageRepo_MarkMessageAsProcessed(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	message := entity.Message{
		Message: "test message",
	}

	id, err := repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	err = repo.MarkMessageAsProcessed(ctx, id)
	require.NoError(t, err)

	fetchedMessage, err := repo.GetMessageById(ctx, id)
	require.NoError(t, err)
	assert.True(t, fetchedMessage.Processed)
}

func TestMessageRepo_GetProcessedMessagesStats(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	_, err := repo.CreateMessage(ctx, entity.Message{Message: "test message 1"})
	require.NoError(t, err)

	_, err = repo.CreateMessage(ctx, entity.Message{Message: "test message 2"})
	require.NoError(t, err)

	id, err := repo.CreateMessage(ctx, entity.Message{Message: "test message 3"})
	require.NoError(t, err)

	err = repo.MarkMessageAsProcessed(ctx, id)
	require.NoError(t, err)

	count, err := repo.GetProcessedMessagesStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)
}

func TestMessageRepo_GetMessageByContent(t *testing.T) {
	teardown := setupPostgres(t)
	defer teardown()

	repo := pgdb.NewMessageRepo(testDB)
	ctx := context.Background()

	message := entity.Message{
		Message: "unique test message",
	}

	id, err := repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	fetchedMessage, err := repo.GetMessageByContent(ctx, message.Message)
	require.NoError(t, err)
	assert.Equal(t, message.Message, fetchedMessage.Message)
	assert.Equal(t, id, fetchedMessage.ID)
}
