package usecase

import (
	"github.com/sergeyptv/post_service/post/internal/ports"
	"log/slog"
)

type post struct {
	log            *slog.Logger
	postRepository ports.PostRepository
}

func NewPostUsecase(log *slog.Logger, postRepository ports.PostRepository) *post {
	return &post{
		log:            log,
		postRepository: postRepository,
	}
}
