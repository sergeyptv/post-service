package http

import (
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"strings"
)

var (
	errDtoInvalid = errors.New("dto is invalid")
)

type userDtoRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func userDtoRegisterToDomain(userDto userDtoRegister) domain.User {
	return domain.User{
		Username: userDto.Username,
		Email:    userDto.Email,
	}
}

// TODO: add validation libruary
func (u *userDtoRegister) Validate() error {
	if strings.TrimSpace(u.Username) == "" ||
		strings.TrimSpace(u.Email) == "" || !strings.Contains(u.Email, "@") ||
		strings.TrimSpace(u.Password) == "" || len(u.Password) < 8 {
		return errDtoInvalid
	}

	return nil
}

type userDtoLogin struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *userDtoLogin) Validate() error {
	if strings.TrimSpace(u.Email) == "" || !strings.Contains(u.Email, "@") ||
		strings.TrimSpace(u.Password) == "" || len(u.Password) < 8 {
		return errDtoInvalid
	}

	return nil
}
