package http

import "net/http"

type router struct {
	Mux *http.ServeMux
}

func NewRouter(h *handler) *router {
	mux := http.NewServeMux()
	mux.Handle("POST /api/register", wrapRateLimiter(h, h.Register))
	mux.Handle("POST /api/login", wrapRateLimiter(h, h.Login))
	mux.Handle("POST /api/logout", http.HandlerFunc(h.Logout))
	mux.Handle("POST /api/refresh", http.HandlerFunc(h.Refresh))

	return &router{
		Mux: mux,
	}
}

func wrapRateLimiter(h *handler, next http.HandlerFunc) http.Handler {
	return h.RateLimiterMiddleware(next)
}
