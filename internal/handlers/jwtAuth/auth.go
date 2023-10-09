package jwtAuth

import "github.com/Genry72/gophermart/internal/models"

type Auther interface {
	GetToken(user *models.User) (string, error)
	ValidateAndParseToken(token string) (int64, string, error)
}
