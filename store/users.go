package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"

	_ "github.com/lib/pq"
)

type UsersStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UsersStore {
	return &UsersStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type User struct {
	ID                 uuid.UUID `db:"id"`
	Email              string    `db:"email"`
	PasswordHashBase64 string    `db:"password_hash"`
	CreatedAt          time.Time `db:"created_at"`
}

func (u *User) CheckPasswordValid(password string) error {
	passwordHash, err := base64.StdEncoding.DecodeString(u.PasswordHashBase64)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	const query = "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;"

	var user User
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash the password: %w", err)
	}
	passwordHashB64 := base64.StdEncoding.EncodeToString(bytes)

	if err := s.db.GetContext(ctx, &user, query, email, passwordHashB64); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *UsersStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	const query = "SELECT * FROM users WHERE email = $1;"

	var user User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (s *UsersStore) FindByID(ctx context.Context, userId uuid.UUID) (*User, error) {
	const query = "SELECT * FROM users WHERE id = $1;"

	var user User
	if err := s.db.GetContext(ctx, &user, query, userId); err != nil {
		return nil, fmt.Errorf("failed to get user by ID(%s): %w", userId, err)
	}

	return &user, nil
}
