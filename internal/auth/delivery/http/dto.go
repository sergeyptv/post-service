package http

import "github.com/sergeyptv/post_service/internal/auth/domain"

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
