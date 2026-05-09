package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type ctxKey string

const userKey ctxKey = "user"

func (h *handler) TokenCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.SplitN(authHeader, " ", 2)

		if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "Bearer") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeaderParts[1]

		user, err := h.jwtParser.Parse(r.Context(), tokenStr)
		if err != nil {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
