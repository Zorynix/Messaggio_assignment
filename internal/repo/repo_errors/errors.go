package repoerrs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInsertFailed  = errors.New("failed to insert record")
	ErrUpdateFailed  = errors.New("failed to update record")
)
