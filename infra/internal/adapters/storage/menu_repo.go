package storage

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/mapper"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"gorm.io/gorm"
)

type menuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) *menuRepo {
	return &menuRepo{db: db}
}

func (r *menuRepo) create(ctx context.Context, menu *model.Menu) error {
	return r.db.Table("menus").WithContext(ctx).Create(menu).Error
}

func (r *menuRepo) GetByID(ctx context.Context, id uint) (*model.Menu, error) {
	var menu types.Menu
	err := r.db.Table("menus").WithContext(ctx).Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}

	return mapper.MenuStorageToDomain(menu), nil
}

func (r *menuRepo) Update(ctx context.Context, menu *model.Menu) error {
	return r.db.Table("menus").WithContext(ctx).Where("id = ?", menu.ID).Updates(menu).Error
}

