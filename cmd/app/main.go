package main

import (
	"context"
	"log"
	"net/http"
	"todo_list/config"
	"todo_list/internal/infrastructure/adapters/logger"
	"todo_list/internal/interfaces/http/server"

	"os"
	"os/signal"

	_ "net/http/pprof"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := logger.New(cfg)

	if cfg.PPROF {
		go func() {
			logger.Info("Starting pprof server on :6060")
			if err := http.ListenAndServe(":6060", nil); err != nil {
				log.Fatalf("pprof server failed: %v", err)
			}
		}()
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataBase, err := server.NewDatabase(ctx, cfg)
	if err != nil {
		logger.Error("Failed to connect to database", err)

		return
	}
	defer dataBase.Close()

	var srv *http.Server

	go func() {
		srv, err = server.Run(ctx, dataBase.Pool, cfg)
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	<-shutdown

	logger.Info("Shutting down server")

	cancel()

	err = server.Close(ctx, srv)
	if err != nil {
		logger.Error("Failed to close server", err)

		return
	}

	dataBase.Close()

	logger.Info("Server shutdown complete")
}
