package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
)

type BranchRepo interface {
	Create(ctx context.Context, branch model.Branch) (uint, error)
	GetByFilter(ctx context.Context, filter *model.BranchFilter) (*model.Branch, error)
}
