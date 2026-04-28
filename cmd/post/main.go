package main

import (
	"context"
	"github.com/sergeyptv/post_service/internal/platform/httpserver"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/post/config"
	postHttp "github.com/sergeyptv/post_service/internal/post/delivery/http"
	"github.com/sergeyptv/post_service/internal/post/repository/postgres"
	"github.com/sergeyptv/post_service/internal/post/usecase"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.App.Env)

	if err := appRun(log, cfg); err != nil {
		panic(err)
	}
}

func appRun(log *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := platformPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	postgresPostRepository := postgres.NewPostgresPostRepository(pool.Db)

	postService := usecase.NewPostService(log, cfg, postgresPostRepository)

	handler := postHttp.NewHandler(postService)

	postServer := httpserver.New(handler, cfg.Server)
}
