package main

import (
	"sync"

	"github.com/orainmers/golangStudy/internal/logger"
	"github.com/orainmers/golangStudy/internal/server"
)

func main() {
	lg := logger.New()

	httpServer := server.New(lg, ":8080")

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
