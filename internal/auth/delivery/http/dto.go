package http

import (
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

var (
	errDtoInvalid = errors.New("dto is invalid")
)

type userDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func userDtoToDomain(userDto userDto) domain.User {
	return domain.User{
		Username: userDto.Username,
		Email:    userDto.Email,
	}
}

func (u *userDto) Validate() error {
	if u.Email == "" || u.Password == "" {
		return errDtoInvalid
	}

	return nil
}
