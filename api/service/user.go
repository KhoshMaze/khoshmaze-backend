package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/jwt"
	timeutils "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/time"
	notifDomain "github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	notifPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/port"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"
	userPort "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/port"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongOTP            = errors.New("wrong otp")
	ErrWrongOTPType        = errors.New("wrong otp type. [0 for register & 1 for login]")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserOnCreate        = errors.New("couldn't create the user")
)

type UserService struct {
	svc           userPort.Service
	notifSvc      notifPort.Service
	authSecret    string
	refreshSecret string
	expMin        uint
	refreshExpMin uint
}

func NewUserService(svc userPort.Service, authSecret, refreshSecret string, expMin, refreshExpMin uint, notifSvc notifPort.Service) *UserService {
	return &UserService{
		svc:           svc,
		notifSvc:      notifSvc,
		authSecret:    authSecret,
		refreshSecret: refreshSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
	}
}

func (s *UserService) SignUp(ctx context.Context, req *pb.UserSignUpRequest, userIP string) (*pb.UserTokenResponse, error) {
	ok, err := s.notifSvc.CheckUserNotifValue(ctx, model.Phone(req.GetPhone()), req.GetOtp())
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrWrongOTP
	}

	if user, _ := s.svc.GetUserByFilter(ctx, &model.UserFilter{
		Phone: req.GetPhone(),
	}); user != nil {
		return nil, ErrUserAlreadyExists
	}
	userId, err := s.svc.CreateUser(ctx, model.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     model.Phone(req.GetPhone()),
	})

	if err != nil {
		return nil, ErrUserOnCreate
	}

	return s.generateTokenResponse(&jwt.UserClaims{UserID: uint(userId), Phone: req.GetPhone(), IP: userIP}, true)

}

func (s *UserService) Login(ctx context.Context, req *pb.UserLoginRequest, userIP string) (*pb.UserTokenResponse, error) {
	ok, err := s.notifSvc.CheckUserNotifValue(ctx, model.Phone(req.GetPhone()), req.GetOtp())
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrWrongOTP
	}

	user, err := s.svc.GetUserByFilter(ctx, &model.UserFilter{
		Phone: req.GetPhone(),
	})

	return s.generateTokenResponse(&jwt.UserClaims{UserID: uint(user.ID), Phone: string(user.Phone), IP: userIP}, true)

}

func (s *UserService) SendOTP(ctx context.Context, req *pb.OtpRequest) error {
	var (
		phone   = req.GetPhone()
		otpType = req.GetType()
	)
	user, err := s.svc.GetUserByFilter(ctx, &model.UserFilter{
		Phone: phone,
	})

	switch otpType {
	case 1: // for register
		if user != nil {
			return ErrUserAlreadyExists
		}
	case 2: // for login
		if err != nil {
			return err
		}
	case 0:
		return ErrWrongOTPType
	}

	code := rand.IntN(999999) + 100000
	return s.notifSvc.Send(ctx, notifDomain.NewNotification(0, fmt.Sprint(code), notifDomain.NotifTypeSMS, true, time.Second*150, model.Phone(phone)))
}

func (s *UserService) Logout(ctx context.Context, token string) error {
	userClaims, err := utils.UserClaimsFromCookies(token, []byte(s.refreshSecret))

	if err != nil {
		return err
	}

	if ok := s.svc.IsBannedToken(ctx, token); ok {
		return ErrInvalidRefreshToken
	}

	err = s.svc.CreateBannedToken(ctx, model.TokenBlacklist{
		Value:     token,
		ExpiresAt: userClaims.ExpiresAt.Time,
		UserID:    model.UserID(userClaims.UserID),
	})

	return err
}

func (s *UserService) RefreshToken(ctx context.Context, token string, userIP string) (*pb.UserTokenResponse, error) {
	userClaims, err := utils.UserClaimsFromCookies(token, []byte(s.refreshSecret))

	if err != nil {
		return nil, err
	}

	if userClaims.IP != userIP {
		return nil, ErrInvalidRefreshToken
	}

	if ok := s.svc.IsBannedToken(ctx, token); ok {
		return nil, ErrInvalidRefreshToken
	}

	if time.Until(userClaims.ExpiresAt.Time)/time.Hour < 12 {

		err = s.svc.CreateBannedToken(ctx, model.TokenBlacklist{
			Value:     token,
			ExpiresAt: userClaims.ExpiresAt.Time,
			UserID:    model.UserID(userClaims.UserID),
		})

		if err != nil {
			return nil, err
		}
		return s.generateTokenResponse(userClaims, true)
	}

	return s.generateTokenResponse(userClaims, false)
}

func (s *UserService) createToken(claims *jwt.UserClaims, isRefresh bool) (string, int64, error) {
	var (
		secret string = s.authSecret
		exp    uint   = s.expMin
	)

	if isRefresh {
		secret = s.refreshSecret
		exp = s.refreshExpMin
	}
	claims.ExpiresAt = jwt5.NewNumericDate(timeutils.AddMinutes(exp, true))
	token, err := jwt.CreateToken([]byte(secret), claims)

	if err != nil {
		return "", 0, err
	}

	return token, int64(exp * 60), nil
}

func (s *UserService) generateTokenResponse(claims *jwt.UserClaims, genRefresh bool) (*pb.UserTokenResponse, error) {

	access, accessMaxAge, err := s.createToken(claims, false)
	if err != nil {
		return nil, err
	}

	var (
		refresh       string
		refreshMaxAge int64
	)
	if genRefresh {

		refresh, refreshMaxAge, err = s.createToken(claims, true)

		if err != nil {
			return nil, err
		}
	}

	return &pb.UserTokenResponse{
		AccessToken:   access,
		RefreshToken:  refresh,
		AccessMaxAge:  accessMaxAge,
		RefreshMaxAge: refreshMaxAge,
	}, nil
}
