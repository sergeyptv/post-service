package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var (
	errMaxLimitAchieved = errors.New("max limit for requests is achieved")
)

func (h *handler) RateLimiterMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("request",
			slog.String("ip", r.URL.Hostname()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var req struct {
			Email string `json:"email"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		err = h.checkRateLimiter(r.Context(), r.URL.Hostname(), h.redisConfig.IpRateLimit, h.redisConfig.IpRateLimiterTtl)
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

		err = h.checkRateLimiter(r.Context(), req.Email, h.redisConfig.EmailRateLimit, h.redisConfig.EmailRateLimiterTtl)
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

		next.ServeHTTP(w, r)
	}
}

func (h *handler) checkRateLimiter(ctx context.Context, limiter string, maxLimit int, ttl time.Duration) error {
	limit, err := h.sessionRepo.Get(ctx, limiter)
	if err != nil {
		if errors.Is(err, repository.ErrDbClientClosed) {
			return err
		}

		_, err = h.sessionRepo.Set(ctx, limiter, "1", ttl)
		if err != nil {
			return err
		}
	} else {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}

		if maxLimit <= limitInt {
			return errMaxLimitAchieved
		}

		limitInt++
		limit = strconv.Itoa(limitInt)

		_, err = h.sessionRepo.Set(ctx, limiter, limit, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}
