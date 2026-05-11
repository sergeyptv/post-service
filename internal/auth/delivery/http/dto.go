package http

import (
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

var (
	errDtoInvalid = errors.New("dto is invalid")
)

type userDtoRegister struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=30"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

func userDtoRegisterToDomain(userDto userDtoRegister) domain.User {
	return domain.User{
		Username: userDto.Username,
		Email:    userDto.Email,
	}
}

type userDtoLogin struct {
	Password string `json:"password" validate:"required,min=8,max=30"`
	Email    string `json:"email" validate:"required,email,max=255"`
}
