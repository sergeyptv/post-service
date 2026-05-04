package http

import (
	"net/http"
)

type router struct {
	Mux *http.ServeMux
}

func NewRouter(handler *handler) *router {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/posts", handler.Create)
	mux.HandleFunc("GET /api/posts/{uuid}", handler.Get)
	mux.HandleFunc("GET /api/posts", handler.List)
	mux.HandleFunc("PATCH /api/posts/{uuid}", handler.Update)
	mux.HandleFunc("DELETE /api/posts/{uuid}", handler.Delete)

	middlewareMux := http.NewServeMux()
	middlewareMux.Handle("/api/", handler.TokenCheckMiddleware(mux))

	return &router{
		Mux: middlewareMux,
	}
}
