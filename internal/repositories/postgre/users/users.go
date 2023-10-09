package users

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserStorage struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewUserStorage(conn *sqlx.DB, log *zap.Logger) *UserStorage {
	return &UserStorage{conn: conn, log: log}
}
