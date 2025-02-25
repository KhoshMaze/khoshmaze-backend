package service

import (
	"context"
	"errors"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/jwt"
	timeutils "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/time"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	svc           userPort.Service
	authSecret    string
	refreshSecret string
	expMin        uint
	refreshExpMin uint
}

func NewUserService(svc userPort.Service, authSecret, refreshSecret string, expMin, refreshExpMin uint) *UserService {
	return &UserService{
		svc:           svc,
		authSecret:    authSecret,
		refreshSecret: refreshSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
	}
}

func (s *UserService) SignUp(ctx context.Context, req *pb.UserSignUpRequest) (*pb.UserTokenResponse, error) {
	userId, err := s.svc.CreateUser(ctx, model.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     model.Phone(req.GetPhone()),
	})

	if err != nil {
		return nil, err
	}

	access, err := s.createToken(uint(userId), req.GetPhone(), false)
	if err != nil {
		return nil, err
	}

	refresh, err := s.createToken(uint(userId), req.GetPhone(), true)

	if err != nil {
		return nil, err
	}

	return &pb.UserTokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}

func (s *UserService) Logout(ctx context.Context, token string) error {
	userClaims, err := utils.UserClaimsFromCookies(token, []byte(s.refreshSecret))

	if err != nil {
		return err
	}

	if ok := s.svc.IsBannedToken(ctx, token); ok {
		return errors.New("invalid refresh token")
	}

	err = s.svc.CreateBannedToken(ctx, model.TokenBlacklist{
		Value:     token,
		ExpiresAt: userClaims.ExpiresAt.Time,
		UserID:    model.UserID(userClaims.UserID),
	})

	return err
}

func (s *UserService) RefreshToken(ctx context.Context, token string) (*pb.UserTokenResponse, error) {
	userClaims, err := utils.UserClaimsFromCookies(token, []byte(s.refreshSecret))

	if err != nil {
		return nil, err
	}

	if ok := s.svc.IsBannedToken(ctx, token); ok {
		return nil, errors.New("invalid refresh token")
	}

	access, err := s.createToken(userClaims.UserID, userClaims.Phone, false)

	if err != nil {
		return nil, err
	}

	err = s.svc.CreateBannedToken(ctx, model.TokenBlacklist{
		Value:     token,
		ExpiresAt: userClaims.ExpiresAt.Time,
		UserID:    model.UserID(userClaims.UserID),
	})

	if err != nil {
		return nil, err
	}

	var refresh string
	if time.Until(userClaims.ExpiresAt.Time) < 12 {
		refresh, err = s.createToken(userClaims.UserID, userClaims.Phone, true)
		if err != nil {
			return nil, err
		}
	}

	return &pb.UserTokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *UserService) createToken(userID uint, phone string, isRefresh bool) (string, error) {
	var (
		secret string = s.authSecret
		exp uint = s.expMin
	)

	if isRefresh {
		secret = s.refreshSecret
		exp = s.refreshExpMin
	}

	token, err := jwt.CreateToken([]byte(secret), &jwt.UserClaims{
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: jwt5.NewNumericDate(timeutils.AddMinutes(exp, true)),
		},
		UserID: uint(userID),
		Phone:  phone,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}
