package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const maxBytes = 10 << 20

var (
	errMaxLimitAchieved = errors.New("max limit for requests is achieved")
)

func (h *handler) RateLimiterMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("request",
			slog.String("ip:port", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			var maxBytesErr *http.MaxBytesError

			if errors.As(err, &maxBytesErr) {
				http.Error(w, "Payload too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		var req struct {
			Email string `json:"email"`
		}

		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			http.Error(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		reqIp := strings.Split(r.RemoteAddr, ":")[0]

		err = h.checkRateLimiter(r.Context(), fmt.Sprintf("rl:ip:%s", reqIp), h.redisConfig.IpRateLimit, h.redisConfig.IpRateLimiterTtl)
		if err != nil {
			switch {
			case errors.Is(err, errMaxLimitAchieved):
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return

			case errors.Is(err, repository.ErrDbClientClosed):
				http.Error(w, "bad gateway", http.StatusBadGateway)
				return

			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		err = h.checkRateLimiter(r.Context(), fmt.Sprintf("rl:email:%s", req.Email), h.redisConfig.EmailRateLimit, h.redisConfig.EmailRateLimiterTtl)
		if err != nil {
			switch {
			case errors.Is(err, errMaxLimitAchieved):
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return

			case errors.Is(err, repository.ErrDbClientClosed):
				http.Error(w, "bad gateway", http.StatusBadGateway)
				return

			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		r.ContentLength = int64(len(bodyBytes))

		next.ServeHTTP(w, r)
	}
}

func (h *handler) checkRateLimiter(ctx context.Context, limiter string, maxLimit int, ttl time.Duration) error {
	limit, err := h.rateLimitRepo.Increment(ctx, limiter, ttl)
	if err != nil {
		return err
	}

	if limit >= maxLimit {
		return errMaxLimitAchieved
	}

	return nil
}
