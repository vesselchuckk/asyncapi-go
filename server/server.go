package server

import (
	"context"
	"github.com/astroniumm/go-asyncapi/config"
	"github.com/astroniumm/go-asyncapi/store"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Config *config.Config
	Logger *slog.Logger
	Store  *store.Store
}

func NewServer(config *config.Config, logger *slog.Logger, store *store.Store) *Server {
	return &Server{
		Config: config,
		Logger: logger,
		Store:  store,
	}
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, "failed to ping", http.StatusBadRequest)
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", s.Ping)
	mux.HandleFunc("POST /auth/signup", s.SignUpHandler())

	middleware := NewLoggerMiddleware(s.Logger)
	server := &http.Server{
		Addr:    net.JoinHostPort(s.Config.ServerHost, s.Config.ServerPort),
		Handler: middleware(mux),
	}

	go func() {
		s.Logger.Info("server is running", "port", s.Config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("server refused to listen and serve")
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.Logger.Error("apiserver failed to shutdown")
		}
	}()
	wg.Wait()

	return server.ListenAndServe()
}
