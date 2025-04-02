package port

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
)

type MenuRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Menu, error)
	Update(ctx context.Context, menu *model.Menu) error
}
