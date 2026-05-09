package http

import (
	"encoding/json"
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/ports"
	"github.com/sergeyptv/post_service/internal/platform/redis"
	"log/slog"
	"net/http"
)

type handler struct {
	log           *slog.Logger
	redisConfig   redis.Config
	rateLimitRepo ports.RateLimitRepository
	usecase       ports.Usecase
}

func NewHandler(log *slog.Logger, redisConfig redis.Config, rateLimitRepo ports.RateLimitRepository, usecase ports.Usecase) *handler {
	return &handler{
		log:           log,
		redisConfig:   redisConfig,
		rateLimitRepo: rateLimitRepo,
		usecase:       usecase,
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var user userDtoRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "can not validate request", http.StatusBadRequest)
		return
	}

	userUuid, err := h.usecase.Register(r.Context(), userDtoRegisterToDomain(user), user.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, "user with this username or email already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	res := struct {
		UserUuid string `json:"user_uuid"`
	}{
		UserUuid: userUuid,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resBytes)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var user userDtoLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, "can not validate request", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.usecase.Login(r.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, "email or password is invalid", http.StatusBadRequest)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Secure:   false,
			HttpOnly: false,
			SameSite: 0,
		})

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Cannot get cookie", http.StatusUnauthorized)
		return
	}

	err = cookie.Valid()
	if err != nil {
		http.Error(w, "Cookie is invalid", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.usecase.Refresh(r.Context(), cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenInvalid) ||
			errors.Is(err, domain.ErrIssIncorrect) ||
			errors.Is(err, domain.ErrKidNotSet) ||
			errors.Is(err, domain.ErrKidIncorrect) ||
			errors.Is(err, domain.ErrExpFired):
			http.Error(w, "token is invalid", http.StatusForbidden)
			return

		case errors.Is(err, domain.ErrClientNotRespond):
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Secure:   false,
			HttpOnly: false,
			SameSite: 0,
		})

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Cannot get cookie", http.StatusUnauthorized)
		return
	}

	err = cookie.Valid()
	if err != nil {
		http.Error(w, "Cookie is invalid", http.StatusUnauthorized)
		return
	}

	err = h.usecase.Logout(r.Context(), cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenInvalid) ||
			errors.Is(err, domain.ErrIssIncorrect) ||
			errors.Is(err, domain.ErrKidNotSet) ||
			errors.Is(err, domain.ErrKidIncorrect) ||
			errors.Is(err, domain.ErrExpFired):
			http.Error(w, "token is invalid", http.StatusForbidden)
			return

		case errors.Is(err, domain.ErrClientNotRespond):
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
