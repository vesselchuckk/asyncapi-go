package main

import (
	"context"
	"github.com/astroniumm/go-asyncapi/config"
	"github.com/astroniumm/go-asyncapi/server"
	"github.com/astroniumm/go-asyncapi/store"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	conf, err := config.New()
	if err != nil {
		return nil
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	db, err := store.NewPostgresDB()
	if err != nil {
		return err
	}

	dataStore := store.New(db)
	jwtManager := server.NewJWTManager(conf)
	server := server.NewServer(conf, logger, dataStore, jwtManager)
	if err := server.Run(ctx); err != nil {
		return err
	}
	return nil
}
