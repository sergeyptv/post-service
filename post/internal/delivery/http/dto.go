package http

import (
	"github.com/sergeyptv/post_service/post/internal/domain"
)

type createPost struct {
	Description string `json:"description" validate:"required,min=1,max=3000"`
	Media       string `json:"media"`
}

func createPostToDomain(post createPost) domain.Post {
	return domain.Post{
		Description: post.Description,
		Media:       &post.Media,
	}
}

type updatePost struct {
	Uuid        string `validate:"required"`
	Description string `json:"description" validate:"required_without=Media"`
	Media       string `json:"media"`
}

func updatePostToDomain(post updatePost) domain.Post {
	return domain.Post{
		Uuid:        post.Uuid,
		Description: post.Description,
		Media:       &post.Media,
	}
}
