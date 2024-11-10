package auth

import (
	"context"
	"errors"
	token "github.com/2pizzzza/plumbing/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *Auth) Register(ctx context.Context, username, password string) error {
	return s.authRepository.CreateUser(ctx, username, password)
}

func (s *Auth) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.authRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("incorrect password")
	}

	generateToken, err := token.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return generateToken, nil
}
