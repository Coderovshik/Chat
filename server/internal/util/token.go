package util

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnknownClaimsType = errors.New("unknown claims type, cannot proceed")
)

type UserClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTSignedString(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseUserClaims(ss string, key string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, ErrUnknownClaimsType
	}
}
