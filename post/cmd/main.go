package main

import (
	"context"
	"errors"
	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
	grpcClient "github.com/sergeyptv/post_service/platform/grpc_client"
	httpServer "github.com/sergeyptv/post_service/platform/http_server"
	"github.com/sergeyptv/post_service/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/platform/postgres"
	authGrpcClient "github.com/sergeyptv/post_service/post/internal/auth/client"
	"github.com/sergeyptv/post_service/post/internal/auth/jwt"
	"github.com/sergeyptv/post_service/post/internal/config"
	http3 "github.com/sergeyptv/post_service/post/internal/delivery/http"
	"github.com/sergeyptv/post_service/post/internal/repository/postgres"
	"github.com/sergeyptv/post_service/post/internal/usecase"
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

	jwtParser := jwt.NewJwtTokenParser(log, cfg.Jwt, cache, authClient)

	postUsecase := usecase.NewPostUsecase(log, postgresPostRepository)

	handler := http3.NewHandler(log, postUsecase, jwtParser)
	router := http3.NewRouter(handler)

	postServer := httpServer.New(router.Mux, cfg.Server)
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
