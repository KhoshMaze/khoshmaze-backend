package service

import (
	"context"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/common"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	restaurantPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/port"
)

type RestaurantService struct {
	svc restaurantPort.Service
}

func NewRestaurantService(svc restaurantPort.Service) *RestaurantService {
	return &RestaurantService{svc: svc}
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, ownerID uint, req *pb.CreateRestaurantRequest) (uint, error) {

	restaurant := model.Restaurant{
		Name:    req.GetName(),
		URL:     req.GetUrl(),
		OwnerID: ownerID,
	}
	return s.svc.CreateRestaurant(ctx, restaurant)
}

func (s *RestaurantService) GetBranch(ctx context.Context, restaurant string, id uint) (*pb.Branch, error) {
	branch, err := s.svc.GetBranchByFilter(ctx, &model.BranchFilter{
		ID:             id,
		RestaurantName: restaurant,
	})

	if err != nil || branch == nil {
		return nil, err
	}
	return &pb.Branch{
		Id:             int64(branch.ID),
		Address:        branch.Address,
		Phone:          branch.Phone,
		Restaurant:     restaurant,
		PrimaryColor:   branch.PrimaryColor,
		SecondaryColor: branch.SecondaryColor,
	}, nil
}

func (s *RestaurantService) CreateBranch(ctx context.Context, req *pb.CreateBranchRequest) (uint, error) {
	branch := model.Branch{
		RestaurantID: uint(req.GetRestaurantID()),
		Address:      req.GetAddress(),
		Phone:        req.GetPhone(),
	}
	return s.svc.CreateBranch(ctx, branch)
}

func (s *RestaurantService) GetAllRestaurants(ctx context.Context, page, pageSize int) (*pb.GetAllRestaurantsResponse, error) {
	pagination := common.NewPagination(page, pageSize)

	result, err := s.svc.GetAllRestaurants(ctx, pagination)
	if err != nil {
		return nil, err
	}

	restaurants := make([]*pb.Restaurant, len(result.Items))
	for i, r := range result.Items {
		restaurants[i] = &pb.Restaurant{
			Id:      int64(r.ID),
			Name:    r.Name,
			Url:     r.URL,
			OwnerId: int64(r.OwnerID),
		}
	}

	return &pb.GetAllRestaurantsResponse{
		Restaurants: restaurants,
		PaginationInfo: &pb.Pagination{
			Page:       int32(result.Page),
			PageSize:   int32(result.PageSize),
			TotalItems: result.TotalItems,
			TotalPages: int32(result.TotalPages),
		},
	}, nil
}
