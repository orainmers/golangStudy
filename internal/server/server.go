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

	app app
}

func New(lg *slog.Logger, addr string, app app) *Server {
	s := Server{
		lg:  lg.WithGroup(pkgName),
		app: app,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/time", s.getTimeHandler)

			r.Route("/person", func(r chi.Router) {
				r.Post("/add", s.addPersonHandler)
			})
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
	s.lg.Info("http server is running...")
	return s.server.ListenAndServe()
}
