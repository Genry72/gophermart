package jwttoken

import (
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int64
	UserName string
}

type JwtToken struct {
	tokenKey string
	lifetime time.Duration
}

func NewJwtToken(tokenKey string, lifetime time.Duration) *JwtToken {
	return &JwtToken{tokenKey: tokenKey, lifetime: lifetime}
}

func (j *JwtToken) GetToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.lifetime)),
		},

		UserID:   user.UserID,
		UserName: user.Username,
	})

	tokenString, err := token.SignedString([]byte(j.tokenKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return tokenString, nil
}

func (j *JwtToken) ValidateAndParseToken(token string) (int64, string, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.tokenKey), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("jwt.ParseWithClaims: %w", err)
	}

	return claims.UserID, claims.UserName, nil
}
