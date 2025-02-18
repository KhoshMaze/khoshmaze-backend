package jwt

import (
	"errors"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

const UserClaimkey = "User-Claims"

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt5.NewWithClaims(jwt5.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenStr string, secret []byte) (*UserClaims, error) {
	token, err := jwt5.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt5.Token) (interface{}, error) {
		return secret, nil
	})

	if token == nil {
		return nil, errors.New("invalid token")
	}

	var claim *UserClaims

	if token.Claims != nil {
		cc, ok := token.Claims.(*UserClaims)
		if ok {
			claim = cc
		}
	}

	if err != nil {
		return claim, err
	}

	if !token.Valid {
		return claim, err
	}

	return claim, nil
}
