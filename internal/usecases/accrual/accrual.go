package accrual

import (
	"context"
	"github.com/Genry72/gophermart/internal/repositories"
	"github.com/Genry72/gophermart/internal/repositories/accrual"
	"github.com/Genry72/gophermart/internal/repositories/postgre"
	"go.uber.org/zap"
	"time"
)

type Accrual struct {
	log     *zap.Logger
	repo    repositories.Accrualer
	doneCtx context.Context // Контекст указывающий на прекращение работы
}

func NewAccrual(repo *postgre.PGStorage, accuralHost string, limitConcurentRequest int, log *zap.Logger) *Accrual {
	return &Accrual{
		log:  log,
		repo: accrual.NewAccrualsStorage(repo.Conn, accuralHost, limitConcurentRequest, log),
	}
}

// Start запуск обновления статусов
func (o *Accrual) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	donectx, cancel := context.WithCancel(ctx)

	o.doneCtx = donectx

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()

			o.log.Info("accrutal stopped")

			cancel()

			return

		case <-ticker.C:
			ids, err := o.repo.GetUnprocessedOrders(ctx)
			if err != nil {
				o.log.Error("o.getUnprocessedOrders", zap.Error(err))
				continue
			}

			if err := o.repo.WriteStatus(ctx, o.repo.GetAccrualInfo(ctx, ids)); err != nil {
				o.log.Error("o.writeStatus(ctx, o.getAccrualInfo(ctx, ids))", zap.Error(err))
			}
		}
	}
}

// WaitDone ожидание окончания всех задач по обновлению статусов.
// Новые заказы не берутся в работу по обновлению при отмене контексата
func (o *Accrual) WaitDone(ctx context.Context) {
	select {
	case <-ctx.Done():
		o.log.Error("не дождались завершения всех задач AccrualsStorage")
	case <-o.doneCtx.Done():
		o.log.Info("AccrualsStorage succes stoped")
	}
}
