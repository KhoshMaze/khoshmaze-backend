package user

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user model.User) (model.UserID, error) {
	if err := user.Validate(); err != nil {
		return 0, errors.New("couldn't create user")
	}
	userID, err := s.repo.Create(ctx, user)
	return userID, err
}

func (s *service) IsBannedToken(ctx context.Context, value string) bool {
	err := s.repo.IsBannedToken(ctx, value)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (s *service) CreateBannedToken(ctx context.Context, token model.TokenBlacklist) error {
	return s.repo.CreateBannedToken(ctx, token)
}
