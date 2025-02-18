package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ServerResponse[T any] struct {
	Data    *T     `json:"data"`
	Message string `json:"message,omitempty"`
}

func (r SignUpRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required to sign up")
	}
	if r.Password == "" {
		return errors.New("password is required to sign up")
	}
	return nil
}

func (s *Server) SignUpHandler() http.HandlerFunc {
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		
		req, err := decode[SignUpRequest](r)
		if err != nil {
			return NewErrWithStatus(http.StatusBadRequest, err)
		}

		existingUser, err := s.Store.Users.FindByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		}
		if existingUser != nil {
			return NewErrWithStatus(http.StatusConflict, fmt.Errorf("user already exists: %v", err))
		}

		_, err = s.Store.Users.CreateUser(r.Context(), req.Email, req.Password)
		if err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, fmt.Errorf("failed to create user in database: %v", err))
		}

		if err := encode[ServerResponse[struct{}]](ServerResponse[struct{}]{
			Message: "successfully signed up user",
		}, http.StatusCreated, w); err != nil {
			return NewErrWithStatus(http.StatusInternalServerError, err)
		}

		return nil
	})
}
