package accrual

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type AccrualsStorage struct {
	conn                  *sqlx.DB
	log                   *zap.Logger
	httpClient            *resty.Client
	limitConcurentRequest chan struct{} // Количество одновременных запросов в accrual

}

const accrualURL = "/api/orders"

func NewAccrualsStorage(conn *sqlx.DB, accuralHost string, limitConcurentRequest int, log *zap.Logger) *AccrualsStorage {
	restyClient := resty.New()

	restyClient.SetBaseURL(accuralHost + accrualURL)

	restyClient.SetTimeout(time.Second)

	restyClient.SetRetryWaitTime(2 * time.Second)

	restyClient.SetRetryCount(5)

	if limitConcurentRequest == 0 {
		limitConcurentRequest = 1
	}

	return &AccrualsStorage{
		conn:                  conn,
		log:                   log,
		httpClient:            restyClient,
		limitConcurentRequest: make(chan struct{}, limitConcurentRequest),
	}
}

// GetUnprocessedOrders id заказов по которым нужны проверки статусов
func (o *AccrualsStorage) GetUnprocessedOrders(ctx context.Context) ([]int64, error) {
	query := "select order_id from orders where status in ($1, $2,$3)"

	ids := make([]int64, 0)

	err := o.conn.SelectContext(ctx, &ids, query, models.OrderStatusNew, models.OrderStatusRegistered, models.OrderStatusProcessing)
	if err != nil {
		return nil, fmt.Errorf("o.conn.SelectContext: %w", err)
	}

	return ids, nil
}

// GetAccrualInfo получение информации по заказам из accrual
func (o *AccrualsStorage) GetAccrualInfo(ctx context.Context, orderIDs []int64) models.ResponseAccruals {
	result := make([]*models.ResponseAccrual, len(orderIDs))

	wg := sync.WaitGroup{}

	for i := range orderIDs {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return

			case o.limitConcurentRequest <- struct{}{}:
				defer func() {
					<-o.limitConcurentRequest
				}()

				accuralStatus := &models.ResponseAccrual{}

				resp, err := o.httpClient.R().SetContext(ctx).SetResult(accuralStatus).Get(fmt.Sprint(orderIDs[i]))
				if err != nil {
					o.log.Error("o.httpClient.R().SetContext(ctx).SetResult(result).Get", zap.Error(err))
					return
				}

				if resp.StatusCode() != http.StatusOK {
					o.log.Error("resp.StatusCode", zap.Error(myerrors.ErrStatusCodeNotCorrect), zap.Int("code", resp.StatusCode()))
					return
				}

				result[i] = accuralStatus
			}
		}(i)
	}

	wg.Wait()
	return result
}

func (o *AccrualsStorage) WriteStatus(ctx context.Context, src models.ResponseAccruals) error {
	if len(src) == 0 {
		return nil
	}

	tx, err := o.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.conn.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			o.log.Error("tx.Rollback", zap.Error(err))
		}
	}()

	var orderIDs pq.StringArray = src.GetOrderIDs()

	if len(orderIDs) == 0 {
		return nil
	}

	// Получаем текущие балансы пользователей и блокируем на запись строки с их записями

	lockBalanceQuery := `
SELECT user_id,
       accruals,
       drawal,
       current_balance,
       last_update 
FROM user_balance  
         WHERE user_id in (
         select user_id from orders 
                        where order_id in ($1)
                        ) FOR UPDATE;`

	// Ключ id пользователя, значение - его баланс
	currentBalances := make(map[int64]models.UserBalance, len(orderIDs))

	rowsBalance, err := tx.QueryxContext(ctx, lockBalanceQuery, orderIDs)
	if err != nil {
		fmt.Println(lockBalanceQuery)
		return fmt.Errorf("lock user_balance: :%w", err)
	}

	defer func() {
		if err := rowsBalance.Close(); err != nil {
			o.log.Error("rowsBalance.Close", zap.Error(err))
		}
	}()

	for rowsBalance.Next() {
		var balance models.UserBalance
		if err := rowsBalance.StructScan(&balance); err != nil {
			return fmt.Errorf("rowsBalance.StructScan: %w", err)
		}
		currentBalances[balance.UserID] = balance
	}

	if err := rowsBalance.Err(); err != nil {
		return fmt.Errorf("rowsBalance.Err: %w", err)
	}

	// Обновляем заказы, обновляем мапу с балансом пользователя

	updateOrdersQuery := `
UPDATE orders set status=$2, accrual=$3, updated_at=now()
where order_id = $1 returning user_id`

	for i := range src {
		if src[i] == nil {
			continue
		}

		rowOrder := tx.QueryRowContext(ctx, updateOrdersQuery, src[i].OrderID, src[i].Status, src[i].Accrual)

		var userID int64

		if src[i].Accrual == 0 {
			continue
		}

		if err := rowOrder.Scan(&userID); err != nil {
			return fmt.Errorf("rowOrder.Scan: %w", err)
		}

		balance := currentBalances[userID]

		balance.Accruals += src[i].Accrual

		balance.CurrentBalance += src[i].Accrual

		currentBalances[userID] = balance
	}

	// Пишем в базу обновленный баланс пользователя

	addBalanceQuery := `
update user_balance set accruals = $1, current_balance =$2, last_update = now() where user_id = $3
`

	for uID := range currentBalances {
		if _, err := tx.ExecContext(ctx,
			addBalanceQuery, currentBalances[uID].Accruals, currentBalances[uID].CurrentBalance, uID); err != nil {
			return fmt.Errorf("addBalance: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
