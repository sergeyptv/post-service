package http

import "net/http"

type router struct {
	Mux *http.ServeMux
}

func NewRouter(h *handler) *router {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", h.RateLimiterMiddleware(http.HandlerFunc(h.Register)))
	mux.HandleFunc("POST /api/login", h.RateLimiterMiddleware(http.HandlerFunc(h.Login)))
	mux.HandleFunc("POST /api/logout", h.Logout)
	mux.HandleFunc("POST /api/refresh", h.Refresh)

	return &router{
		Mux: mux,
	}
}
