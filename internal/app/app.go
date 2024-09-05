package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/orainmers/golangStudy/internal/models"
)

const moduleName = "app"

type store interface {
	AddPerson(person *models.Person) error
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

func (a *App) CreatePerson(person *models.Person) error {
	person.ID = uuid.New()

	t := time.Now()
	person.CreatedAt = t
	person.UpdatedAt = t
	person.IsDeleted = false

	if err := a.store.AddPerson(person); err != nil {
		return fmt.Errorf("a.store.AddPerson(...): %v", err)
	}

	return nil
}
