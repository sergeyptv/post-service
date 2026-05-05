package http

import "net/http"

type router struct {
	Mux *http.ServeMux
}

func NewRouter(handler *handler) *router {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", handler.Register)
	mux.HandleFunc("POST /api/login", handler.Login)

	return &router{
		Mux: mux,
	}
}
