package users

import (
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"github.com/Genry72/gophermart/internal/repositories/postgre/users"
	"go.uber.org/zap"
)

type Users struct {
	log  *zap.Logger
	repo repositories.Userser
}

func NewUsers(repo *postgre.PGStorage, log *zap.Logger) *Users {
	return &Users{
		log:  log,
		repo: users.NewUserStorage(repo.Conn, log),
	}
}
