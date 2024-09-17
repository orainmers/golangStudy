package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/orainmers/golangStudy/internal/models"
)

const moduleName = "service"

type store interface {
	AddPerson(person *models.Person) (uuid.UUID, error)
}

type App struct {
	lg    *slog.Logger
	store store
}

func New(lg *slog.Logger, store store) *App {
	return &App{
		lg:    lg.With("module", moduleName),
		store: store,
	}
}

func (a *App) CreatePerson(person *models.Person) (uuid.UUID, error) {
	person.ID = uuid.New()

	t := time.Now()
	person.CreatedAt = t
	person.UpdatedAt = t
	person.IsDeleted = false

	id, err := a.store.AddPerson(person)
	if err != nil {
		return uuid.Nil, fmt.Errorf("a.storage.AddPerson(...): %v", err)
	}

	return id, nil
}
