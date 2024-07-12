package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) getTimeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	currentTime := time.Now().Format(time.RFC1123)

	if err := json.NewEncoder(w).Encode(currentTime); err != nil {
		s.lg.Warn("getTimeHandler(...)", "error", err)
	}
}
