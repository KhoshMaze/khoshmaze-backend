package service

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/port"
)

type MenuService struct {
	svc port.Service
}

func NewMenuService(svc port.Service) *MenuService {
	return &MenuService{svc: svc}
}

func (s *MenuService) GetFoods(ctx context.Context, menuID uint, page, pageSize int) (*pb.GetAllFoodsResponse, error) {
	pagination := common.NewPagination(page, pageSize)

	result, err := s.svc.GetAllFoods(ctx, pagination, menuID)
	if err != nil {
		return nil, err
	}
	foods := make([]*pb.Food, len(result.Items))

	for i, r := range result.Items {
		foods[i] = &pb.Food{
			Id:          int64(r.ID),
			Name:        r.Name,
			Description: r.Description,
			Type:        r.Type,
			IsAvailable: r.IsAvailable,
			Price:       r.Price,
		}
	}

	return &pb.GetAllFoodsResponse{
		Extra: &pb.GetAllFoodsResponse_Extra{
			MenuID: int64(menuID),
			// BranchID:   int64(branchID),
			// RestaurantID: int64(restaurantID),
		},
		Foods: foods,
		PaginationInfo: &pb.Pagination{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			TotalItems: result.TotalItems,
			TotalPages: int32(result.TotalPages),
		},
	}, nil
}

func (s *MenuService) GetMenuByID(ctx context.Context, id uint) (*model.Menu, error) {
	return s.GetMenuByID(ctx, id)
}

func (s *MenuService) AddFood(ctx context.Context, req *pb.CreateFoodRequest) (uint, error) {
	food := model.Food{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Type:        req.GetType(),
		Price:       req.GetPrice(),
		MenuID:      uint(req.GetMenuID()),
	}

	if err := food.Validate(); err != nil {
		return 0, err
	}

	return s.svc.AddFoodToMenu(ctx, food)

}

// TODO: add price update checking logic to reduce db unnecessary queries
func (s *MenuService) UpdateFood(ctx context.Context, req *pb.Food) error {
	food := model.Food{
		ID:          uint(req.GetId()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Type:        req.GetType(),
		Price:       req.GetPrice(),
		IsAvailable: req.GetIsAvailable(),
	}

	if err := food.Validate(); err != nil {
		return err
	}

	return s.svc.UpdateFoodInMenu(ctx, food)
}

func (s *MenuService) DeleteFood(ctx context.Context, id uint) error {
	return s.svc.DeleteFoodFromMenu(ctx, id)
}

func (s *MenuService) GetFood(ctx context.Context, id uint) (*pb.Food, error) {
	food, err := s.svc.GetFoodByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &pb.Food{
		Id:          int64(food.ID),
		Name:        food.Name,
		Description: food.Description,
		Type:        food.Type,
		Price:       food.Price,
		IsAvailable: food.IsAvailable,
	}, nil
}
