package main

import (
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/orainmers/golangStudy/internal/app"
	"github.com/orainmers/golangStudy/internal/logger"
	"github.com/orainmers/golangStudy/internal/server"
	"github.com/orainmers/golangStudy/internal/store"
)

const (
	address  = "127.0.0.1:5432"
	username = "postgres"
	password = "postgres"
	database = "postgres"
)

func main() {
	lg := logger.New()

	psql, err := store.New(lg, username, password, address, database)
	if err != nil {
		lg.Error("Failed to connect to database",
			"error", err)
		return
	}

	defer func() {
		if err = psql.Close(); err != nil {
			lg.Error("Failed to close",
				"error", err)
		}
	}()

	if err = psql.DummyMigration(); err != nil {
		lg.Error("Failed to migrate",
			"error", err)
		return
	}

	application := app.New(lg, psql)

	httpServer := server.New(lg, ":8080", application)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := httpServer.Run(); err != nil {
			lg.Error("server", "error", err)
		}
	}()

	lg.Info("Service is running ...")

	wg.Wait()
}
