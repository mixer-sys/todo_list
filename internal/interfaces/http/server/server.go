package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"todo_list/config"

	router "todo_list/internal/infrastructure/adapters/router"
	"todo_list/internal/interfaces/http/middleware"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Run(ctx context.Context,
	dataBase *pgxpool.Pool, cfg *config.Config) (
	*http.Server, error) {
	r := router.NewRouter(dataBase, cfg)

	address := ":" + cfg.ServerPort

	srv := &http.Server{
		Addr:              address,
		Handler:           middleware.TokenMiddleware(cfg, r),
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeoutSecond) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server listen error: ",
				slog.String("error", err.Error()),
				slog.String("address", address))

			return
		}
	}()

	return srv, nil
}

func Close(ctx context.Context,
	srv *http.Server) error {

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error: ",
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, cfg *config.Config) (*Database, error) {

	pool, err := pgxpool.Connect(ctx, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser,
		cfg.DBPassword, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
