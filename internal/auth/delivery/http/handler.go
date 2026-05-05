package http

import (
	"encoding/json"
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/ports"
	"net/http"
)

type handler struct {
	usecase ports.Usecase
}

func NewHandler(usecase ports.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.InputUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	userUuid, err := h.usecase.Register(r.Context(), user)
	if err != nil {
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
	var user domain.InputUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(r.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, "email or password is invalid", http.StatusBadRequest)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resBytes)
}
