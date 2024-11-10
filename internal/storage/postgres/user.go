package postgres

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *DB) CreateUser(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.Pool.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func (s *DB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := s.Pool.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
