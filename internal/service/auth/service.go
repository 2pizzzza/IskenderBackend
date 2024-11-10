package auth

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
)

type Auth struct {
	log            *slog.Logger
	tokenTTL       string
	authRepository AuthReposiotory
}

type AuthReposiotory interface {
	CreateUser(ctx context.Context, username, password string) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

func New(
	log *slog.Logger,
	tokenTTL string,
	authRep AuthReposiotory,
) *Auth {
	return &Auth{
		log:            log,
		tokenTTL:       tokenTTL,
		authRepository: authRep,
	}
}
