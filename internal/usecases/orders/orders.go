package orders

import (
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/postgre/orders"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Orders struct {
	log  *zap.Logger
	repo repositories.Orderer
}

func NewOrders(conn *sqlx.DB, log *zap.Logger) *Orders {
	return &Orders{
		log:  log,
		repo: orders.NewOrderStorage(conn, log),
	}
}
