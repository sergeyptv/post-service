package usecase

import (
	"github.com/sergeyptv/post_service/internal/post/config"
	"github.com/sergeyptv/post_service/internal/post/ports"
	"log/slog"
)

type post struct {
	log            *slog.Logger
	cfg            *config.Config
	postRepository ports.PostRepository
}

func NewPostService(log *slog.Logger, cfg *config.Config, postRepository ports.PostRepository) *post {
	return &post{
		log:            log,
		cfg:            cfg,
		postRepository: postRepository,
	}
}
