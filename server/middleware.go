package server

import (
	"context"
	"github.com/astroniumm/go-asyncapi/store"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
)

func NewLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("http req", "path", r.Method+" "+r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

type userCtxKey struct{}

func WithUser(ctx context.Context, user *store.User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func NewAuthMiddleware(jwtManager *JwtManager, userStore *store.UsersStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}
			// auth header check
			authHeader := r.Header.Get("Authorization")
			var token string
			if parts := strings.Split(authHeader, "Bearer "); len(parts) == 2 {
				token = parts[1]
			}
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			parsedToken, err := jwtManager.Parse(token)
			if err != nil {
				slog.Error("failed to parse token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !jwtManager.IsAccessToken(parsedToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not an access token"))
				return
			}

			userIdstr, err := parsedToken.Claims.GetSubject()
			if err != nil {
				slog.Error("failed to extract user claim from token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userId, err := uuid.Parse(userIdstr)
			if err != nil {
				slog.Error("token subject is invalid", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := userStore.FindByID(r.Context(), userId)
			if err != nil {
				slog.Error("failed to get the user by id", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}
}
