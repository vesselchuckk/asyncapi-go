package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("failed to load .env: %w", err)
	}
	dsn := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to poen db connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
