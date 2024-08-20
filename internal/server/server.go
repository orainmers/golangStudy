package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const pkgName = "server"

type Server struct {
	lg     *slog.Logger
	server *http.Server
}

func New(lg *slog.Logger, addr string) *Server {
	s := Server{
		lg: lg.WithGroup(pkgName),
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/time", s.getTimeHandler)
		})
	})

	s.server = &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &s
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
