package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
	"github.com/sergeyptv/post_service/auth/internal/config"
	"github.com/sergeyptv/post_service/auth/internal/crypto/jwt"
	"github.com/sergeyptv/post_service/auth/internal/delivery/grpc"
	http2 "github.com/sergeyptv/post_service/auth/internal/delivery/http"
	"github.com/sergeyptv/post_service/auth/internal/repository/postgres"
	"github.com/sergeyptv/post_service/auth/internal/repository/redis"
	"github.com/sergeyptv/post_service/auth/internal/usecase"
	"github.com/sergeyptv/post_service/platform/grpc_server"
	"github.com/sergeyptv/post_service/platform/http_server"
	"github.com/sergeyptv/post_service/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/platform/postgres"
	platformRedis "github.com/sergeyptv/post_service/platform/redis"
	"github.com/sergeyptv/post_service/platform/transaction"
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

	jwtTokenSigner := jwt.NewJwtTokenSigner(cfg.Jwt)

	pool, err := platformPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()
	postgresUserRepository := postgres.NewPostgresUserRepository(pool)
	postgresOutboxRepository := postgres.NewPostgresOutboxRepository()

	redisClient, err := platformRedis.New(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer redisClient.Close()
	redisSessionRepository := redis.NewRedisSessionRepository(redisClient)
	redisRateLimitRepository := redis.NewRedisRateLimitRepository(redisClient)

	txWrapper := transaction.New(pool.Db)

	authUsecase := usecase.NewAuthUsecase(log, cfg, postgresUserRepository, postgresOutboxRepository, redisSessionRepository, jwtTokenSigner, txWrapper)

	httpHandler := http2.NewHandler(log, cfg.Redis, redisRateLimitRepository, authUsecase)
	httpRouter := http2.NewRouter(httpHandler)

	authHttpServer := http_server.New(httpRouter.Mux, cfg.HttpServer)
	defer authHttpServer.Close()

	grpcHandler := grpc.NewHandler(authUsecase, cfg.Jwt)

	grpcServer, err := grpc_server.NewServer(cfg.GrpcServer)
	if err != nil {
		return err
	}

	authV1.RegisterAuthServiceServer(grpcServer.Server, grpcHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if gerr := grpcServer.Serve(); gerr != nil {
			log.Error("grpc server failed", logger.Error(gerr))
		}
	}()

	go func() {
		if herr := authHttpServer.Start(); herr != nil && !errors.Is(herr, http.ErrServerClosed) {
			log.Error("http server failed", logger.Error(herr))
		}
	}()

	<-stop

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()

	grpcServer.CloseListener()
	grpcServer.Shutdown()
	return authHttpServer.Shutdown(ctxShutdown)
}
