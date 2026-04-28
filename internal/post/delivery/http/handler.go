package http

import (
	"github.com/sergeyptv/post_service/internal/post/ports"
)

type handler struct {
	usecase ports.Usecase
}

func NewHandler(usecase ports.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}
