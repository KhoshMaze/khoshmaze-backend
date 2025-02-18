package utils

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/jwt"
	"github.com/gofiber/fiber/v2"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func UserClaimsFromLocals(ctx *fiber.Ctx) *jwt.UserClaims {
	if u := ctx.Locals("user"); u != nil {
		userClaims, ok := u.(*jwt5.Token).Claims.(*jwt.UserClaims)
		if ok {
			return userClaims
		}
	}
	return nil
}

func UserClaimsFromCookies(token string, secret []byte) (*jwt.UserClaims, error) {
	if token != "" {
		userclaims, err := jwt.ParseToken(token, secret)
		if err != nil {
			return nil, err
		}
		return userclaims, nil
	}

	return nil, errors.New("empty refresh cookie")
}
