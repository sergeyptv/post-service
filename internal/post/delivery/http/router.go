package http

import (
	"encoding/json"
	"errors"
	"github.com/sergeyptv/post_service/internal/post/domain"
	"net/http"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var post domain.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	if post.Username == "" {
		http.Error(w, "Post cannot have empty username", http.StatusBadRequest)
		return
	}
	if post.Description == "" {
		http.Error(w, "Post cannot have empty description", http.StatusBadRequest)
		return
	}

	postUuid, err := h.usecase.Create(r.Context(), post)
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
	var req struct {
		PostUuid string `json:"post_uuid"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	if req.PostUuid == "" {
		http.Error(w, "Post cannot have empty uuid", http.StatusBadRequest)
		return
	}

	post, err := h.usecase.Get(r.Context(), req.PostUuid)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusBadRequest)
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
	var req struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	postUuids, err := h.usecase.List(r.Context(), req.Username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusBadRequest)
			return

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
	var post domain.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	if post.Uuid == "" {
		http.Error(w, "Post cannot have empty uuid", http.StatusBadRequest)
		return
	}

	if post.Description == "" && len(post.Media) == 0 {
		http.Error(w, "Data for updating is not set", http.StatusBadRequest)
		return
	}

	err = h.usecase.Update(r.Context(), post)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusBadRequest)
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
	var req struct {
		PostUuid string `json:"post_uuid"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	if req.PostUuid == "" {
		http.Error(w, "Post cannot have empty uuid", http.StatusBadRequest)
		return
	}

	err = h.usecase.Delete(r.Context(), req.PostUuid)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPostNotExist):
			http.Error(w, "Post does not exist", http.StatusBadRequest)
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
