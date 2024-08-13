package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/orainmers/golangStudy/internal/models"
	"net/http"
	"time"
)

type DB interface {
	AddPerson(person *models.Person) error
}

func (s *Server) getTimeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	currentTime := time.Now().Format(time.RFC1123)

	if err := json.NewEncoder(w).Encode(currentTime); err != nil {
		s.lg.Warn("getTimeHandler(...)", "error", err)
	}
}

func (s *Server) addPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person models.PersonRequest

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errEncoder := json.NewEncoder(w).Encode(err.Error()); errEncoder != nil {
			s.lg.Warn("addPersonHandler(...)", "error", errEncoder)

			return
		}
	}

	req := models.Person{
		ID:          uuid.New(),
		Name:        person.Name,
		Description: person.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	if err := s.db.AddPerson(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if errEncoder := json.NewEncoder(w).Encode(err.Error()); errEncoder != nil {
			s.lg.Warn("addPersonHandler(...)", "error", errEncoder)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
