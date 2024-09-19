package service

import (
	"fmt"
	"github.com/orainmers/golangStudy/internal/models"
	"log/slog"
)

const moduleName = "service"

// TODO: починить интерфейс
type storage interface {
	AddPerson(person *models.Contact) error
}

type Service struct {
	lg      *slog.Logger
	storage storage
}

func New(lg *slog.Logger, storage storage) *Service {
	return &Service{
		lg:      lg.With("module", moduleName),
		storage: storage,
	}
}

func (a *Service) CreatePerson(person *models.Contact) error {

	person.IsDeleted = false

	// FIXME: сделать проброс id
	if err := a.storage.AddPerson(person); err != nil {
		return fmt.Errorf("a.storage.AddPerson(...): %v", err)
	}

	return nil
}
