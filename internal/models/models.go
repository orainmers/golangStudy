package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID uuid.UUID `db:"id"`

	Name        string `db:"name"`
	Description string `db:"description"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	IsDeleted bool `db:"is_deleted"`
}

type PersonRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
