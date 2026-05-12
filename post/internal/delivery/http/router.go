package http

import (
	"net/http"
)

type router struct {
	Mux *http.ServeMux
}

func NewRouter(h *handler) *router {
	mux := http.NewServeMux()
	mux.Handle("POST /api/posts", wrapAuth(h, h.Create))
	mux.Handle("GET /api/posts/{uuid}", wrapAuth(h, h.Get))
	mux.Handle("GET /api/posts", wrapAuth(h, h.List))
	mux.Handle("PATCH /api/posts/{uuid}", wrapAuth(h, h.Update))
	mux.Handle("DELETE /api/posts/{uuid}", wrapAuth(h, h.Delete))

	return &router{
		Mux: mux,
	}
}

func wrapAuth(h *handler, next http.HandlerFunc) http.Handler {
	return h.TokenCheckMiddleware(next)
}
