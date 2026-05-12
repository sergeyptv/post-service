package http

import (
	"context"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/post/internal/auth/jwt"
	"net/http"
	"strings"
)

type ctxKey string

const userKey ctxKey = "user"

func (h *handler) TokenCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "Bearer") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeaderParts[1]

		user, err := h.jwtParser.Parse(r.Context(), tokenStr, jwt.TypeAccess)
		if err != nil {
			h.log.Warn("Cannot parse token", logger.Error(err))

			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
