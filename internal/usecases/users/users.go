package users

import (
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/postgre/users"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Users struct {
	log  *zap.Logger
	repo repositories.Userser
}

func NewUsers(conn *sqlx.DB, log *zap.Logger) *Users {
	return &Users{
		log:  log,
		repo: users.NewUserStorage(conn, log),
	}
}
