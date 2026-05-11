package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sergeyptv/post_service/internal/post/domain"
	"github.com/sergeyptv/post_service/internal/post/ports"
	"log/slog"
	"net/http"
)

type handler struct {
	log       *slog.Logger
	validate  *validator.Validate
	usecase   ports.Usecase
	jwtParser ports.JwtTokenParser
}

func NewHandler(log *slog.Logger, usecase ports.Usecase, jwtParser ports.JwtTokenParser) *handler {
	return &handler{
		log:       log,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
		usecase:   usecase,
		jwtParser: jwtParser,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var post createPost

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(post)
	if err != nil {
		http.Error(w, "Post cannot have empty description", http.StatusBadRequest)
		return
	}

	userCtx := r.Context().Value(userKey)
	user, ok := userCtx.(domain.User)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	postUuid, err := h.usecase.Create(r.Context(), user, createPostToDomain(post))
	if err != nil {
		if errors.Is(err, domain.ErrBadGateway) {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	res := struct {
		PostUuid string `json:"post_uuid"`
	}{
		PostUuid: postUuid,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resBytes)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	postUuid := r.PathValue("uuid")

	if postUuid == "" {
		http.Error(w, "Post cannot have empty uuid", http.StatusBadRequest)
		return
	}

	userCtx := r.Context().Value(userKey)
	user, ok := userCtx.(domain.User)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	post, err := h.usecase.Get(r.Context(), user, postUuid)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusNotFound)
			return

		case errors.Is(err, domain.ErrBadGateway):
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return

		default:
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}

	resBytes, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(userKey)
	user, ok := userCtx.(domain.User)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	postUuids, err := h.usecase.List(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadGateway):
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return

		default:
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}

	res := struct {
		PostUuids []string `json:"post_uuids"`
	}{
		PostUuids: postUuids,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error preparing the answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var post updatePost

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	post.Uuid = r.PathValue("uuid")

	err = h.validate.Struct(post)
	if err != nil {
		http.Error(w, "Data for updating is not set", http.StatusBadRequest)
		return
	}

	userCtx := r.Context().Value(userKey)
	user, ok := userCtx.(domain.User)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	err = h.usecase.Update(r.Context(), user, updatePostToDomain(post))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusNotFound)
			return

		case errors.Is(err, domain.ErrBadGateway):
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return

		default:
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	postUuid := r.PathValue("uuid")

	if postUuid == "" {
		http.Error(w, "Post cannot have empty uuid", http.StatusBadRequest)
		return
	}

	userCtx := r.Context().Value(userKey)
	user, ok := userCtx.(domain.User)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	err := h.usecase.Delete(r.Context(), user, postUuid)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusNotFound)
			return

		case errors.Is(err, domain.ErrBadGateway):
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return

		default:
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
