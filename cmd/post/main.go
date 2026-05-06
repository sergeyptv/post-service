package main

import (
	"context"
	"errors"
	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
	"github.com/sergeyptv/post_service/internal/platform/grpcClient"
	"github.com/sergeyptv/post_service/internal/platform/httpserver"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	authGrpcClient "github.com/sergeyptv/post_service/internal/post/auth/client"
	"github.com/sergeyptv/post_service/internal/post/auth/jwt"
	"github.com/sergeyptv/post_service/internal/post/config"
	postHttp "github.com/sergeyptv/post_service/internal/post/delivery/http"
	"github.com/sergeyptv/post_service/internal/post/repository/postgres"
	"github.com/sergeyptv/post_service/internal/post/usecase"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	client, err := grpcClient.NewClient(cfg.GrpcClient)
	if err != nil {
		return err
	}
	defer client.Close()
	authServiceClient := authV1.NewAuthServiceClient(client.Conn)
	authClient := authGrpcClient.NewAuthClient(authServiceClient, cfg.Jwt)

	cache := jwt.NewInMemoryCache(cfg.Cache)

	jwtParser := jwt.NewJwtTokenParser(cfg.Jwt, cache, authClient)

	postUsecase := usecase.NewPostUsecase(log, postgresPostRepository)

	handler := postHttp.NewHandler(postUsecase, jwtParser)
	router := postHttp.NewRouter(handler)

	postServer := httpserver.New(router.Mux, cfg.Server)
	defer postServer.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if serr := postServer.Start(); serr != nil && !errors.Is(serr, http.ErrServerClosed) {
			log.Error("server failed", logger.Error(serr))
		}
	}()

	<-stop

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()

	return postServer.Shutdown(ctxShutdown)
}
