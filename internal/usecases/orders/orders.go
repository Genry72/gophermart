package orders

import (
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"github.com/Genry72/gophermart/internal/repositories/postgre/orders"
	"go.uber.org/zap"
)

type Orders struct {
	log  *zap.Logger
	repo repositories.Orderer
}

func NewOrders(repo *postgre.PGStorage, log *zap.Logger) *Orders {
	return &Orders{
		log:  log,
		repo: orders.NewOrderStorage(repo.Conn, log),
	}
}
