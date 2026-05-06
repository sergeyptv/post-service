package main

import (
	"context"
	"errors"
	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
	"github.com/sergeyptv/post_service/internal/auth/config"
	"github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"github.com/sergeyptv/post_service/internal/auth/delivery/grpc"
	authHttp "github.com/sergeyptv/post_service/internal/auth/delivery/http"
	"github.com/sergeyptv/post_service/internal/auth/repository/postgres"
	"github.com/sergeyptv/post_service/internal/auth/usecase"
	"github.com/sergeyptv/post_service/internal/platform/grpcServer"
	"github.com/sergeyptv/post_service/internal/platform/httpserver"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/platform/transaction"
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

	jwtTokenSigner := jwt.NewJwtTokenSigner(cfg.Jwt)

	pool, err := platformPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	postgresUserRepository := postgres.NewPostgresUserRepository(pool)
	postgresTokenRepository := postgres.NewPostgresTokenRepository(pool)
	postgresOutboxRepository := postgres.NewPostgresOutboxRepository(pool)

	txWrapper := transaction.New(pool.Db)

	authUsecase := usecase.NewAuthUsecase(log, postgresUserRepository, postgresOutboxRepository, postgresTokenRepository, jwtTokenSigner, txWrapper)

	httpHandler := authHttp.NewHandler(authUsecase)
	httpRouter := authHttp.NewRouter(httpHandler)

	authHttpServer := httpserver.New(httpRouter.Mux, cfg.HttpServer)
	defer authHttpServer.Close()

	grpcHandler := grpc.NewHandler(authUsecase, cfg.Jwt)

	grpcGrpcServer, err := grpcServer.NewServer(cfg.GrpcServer)
	if err != nil {
		return err
	}

	authV1.RegisterAuthServiceServer(grpcGrpcServer.Server, grpcHandler)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if gerr := grpcGrpcServer.Serve(); gerr != nil {
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

	grpcGrpcServer.Shutdown()
	grpcGrpcServer.CloseListener()
	return authHttpServer.Shutdown(ctxShutdown)
}
