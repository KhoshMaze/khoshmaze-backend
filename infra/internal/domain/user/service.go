package user

import (
	"context"
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	"gorm.io/gorm"
)

var (
	ErrUserOnCreate           = errors.New("error on creating new user")
	ErrUserCreationValidation = errors.New("validation failed")
	ErrUserNotFound           = errors.New("user not found")
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
		return 0, ErrUserOnCreate
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

func (s *service) GetUserByFilter(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	user, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
