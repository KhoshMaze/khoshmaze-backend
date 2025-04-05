package service

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
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
