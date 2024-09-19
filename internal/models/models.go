package models

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID uuid.UUID `db:"id"`

	Name        string `db:"name"`
	Description string `db:"description"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ContactRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
