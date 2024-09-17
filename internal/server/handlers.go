package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/orainmers/golangStudy/internal/models"
)

type service interface {
	CreatePerson(person *models.Person) (uuid.UUID, error)
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

	var req models.PersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errEncoder := json.NewEncoder(w).Encode(err.Error()); errEncoder != nil {
			s.lg.Warn("addPersonHandler(...)", "error", errEncoder)

			return
		}
	}

	person := models.Person{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := s.service.CreatePerson(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if errEncoder := json.NewEncoder(w).Encode(err.Error()); errEncoder != nil {
			s.lg.Warn("addPersonHandler(...)", "error", errEncoder)

			return
		}
	}

	if errEncoder := json.NewEncoder(w).Encode(id); errEncoder != nil {
		s.lg.Warn("addPersonHandler(...)", "error", errEncoder)

		return
	}

	w.WriteHeader(http.StatusOK)
}
