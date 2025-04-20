package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type RefreshTokenStore struct {
	db *sqlx.DB
}

type RefreshToken struct {
	UserID      uuid.UUID `db:"user_id"`
	HashedToken string    `db:"token_hash"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *RefreshTokenStore) getB64Hash(token *jwt.Token) (string, error) {

	h := sha256.New()
	h.Write([]byte(token.Raw))

	hashedBytes := h.Sum(nil)
	base64TokenHash := base64.StdEncoding.EncodeToString(hashedBytes)

	return base64TokenHash, nil
}

func (s *RefreshTokenStore) CreateToken(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const query = `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING *;`

	base64TokenHash, err := s.getB64Hash(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token hash: %w", err)
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to extract expiration time: %w", err)
	}

	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, query, userID, base64TokenHash, expiresAt.Time); err != nil {
		return nil, fmt.Errorf("failed to create refresh JWT token: %w", err)
	}

	return &refreshToken, nil
}

func (s *RefreshTokenStore) ByPrimaryKey(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const query = `SELECT * FROM refresh_tokens WHERE user_id = $1 AND token_hash = $2`

	base64TokenHash, err := s.getB64Hash(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token hash: %w", err)
	}

	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, query, userID, base64TokenHash); err != nil {
		return nil, fmt.Errorf("failed to fetch token for user %s (token: %s) : %w", base64TokenHash, userID, err)
	}

	return &refreshToken, nil
}

func (s *RefreshTokenStore) DeleteUserTokens(ctx context.Context, userId uuid.UUID) (sql.Result, error) {
	const query = `DELETE FROM refresh_tokens WHERE user_id = $1;`
	result, err := s.db.ExecContext(ctx, query, userId)
	if err != nil {
		return result, fmt.Errorf("failed to delete user refresh_tokens: %w", err)
	}

	return result, nil
}
