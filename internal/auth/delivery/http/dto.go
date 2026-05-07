package http

import (
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"strings"
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

// TODO: add prod validator
func (u *userDto) Validate() error {
	if u.Email == "" || u.Password == "" || !strings.Contains(u.Email, "@") {
		return errDtoInvalid
	}

	return nil
}
