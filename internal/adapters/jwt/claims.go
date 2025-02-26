package jwt

import jwt5 "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt5.RegisteredClaims
	UserID  uint
	Phone   string
	IP      string
}
