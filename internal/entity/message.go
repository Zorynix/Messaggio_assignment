package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID          uuid.UUID  `json:"id"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	Processed   bool       `json:"processed"`
	ProcessedAt *time.Time `json:"processed_at"`
}
