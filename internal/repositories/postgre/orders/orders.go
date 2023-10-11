package orders

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type OrderStorage struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewOrderStorage(conn *sqlx.DB, log *zap.Logger) *OrderStorage {
	return &OrderStorage{conn: conn, log: log}
}
