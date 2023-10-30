package balances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/gophermart/internal/models"
	"github.com/Genry72/gophermart/internal/models/myerrors"
	"go.uber.org/zap"
)

// Withdraw Списание средств
func (u *BalanceStorage) Withdraw(ctx context.Context, withdraw *models.Withdraw) error {

	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.conn.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			u.log.Error("tx.Rollback", zap.Error(err))
		}
	}()

	// Получаем текущий баланс пользователя и блокируем на запись строки с его записями
	lockBalanceQuery := `
SELECT user_id,
       accruals,
       drawal,
       current_balance,
       last_update 
FROM user_balance where user_id = $1 FOR UPDATE;`

	var balance models.UserBalance

	rowBalance := tx.QueryRowxContext(ctx, lockBalanceQuery, withdraw.UserID)

	if err := rowBalance.StructScan(&balance); err != nil {
		return fmt.Errorf("rowBalance.StructScan: %w", err)
	}

	if err := rowBalance.Err(); err != nil {
		return fmt.Errorf("rowBalance.Err: %w", err)
	}

	if balance.CurrentBalance-withdraw.Points < 0 {
		return myerrors.ErrNoMoney
	}

	// Пишем в таблицу withdraw
	if _, err := tx.ExecContext(ctx, `
    INSERT INTO withdraw (user_id, order_id, points, date)
    VALUES ($1, $2, $3, NOW())
`, withdraw.UserID, withdraw.Order, withdraw.Points); err != nil {
		return fmt.Errorf("INSERT withdraw: %w", err)
	}

	balance.CurrentBalance -= withdraw.Points

	balance.Drawal += withdraw.Points

	// Пишем в таблицу user_balance
	if _, err := tx.ExecContext(ctx, `
   update user_balance set current_balance = $1, drawal =$2, last_update = now() where user_id = $3
`, balance.CurrentBalance, balance.Drawal, withdraw.UserID); err != nil {
		return fmt.Errorf("addBalance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

// Withdrawals Получение информации о выводе средств
func (u *BalanceStorage) Withdrawals(ctx context.Context, userID int64) ([]*models.Withdraw, error) {
	query := `
select order_id, points, date
from withdraw
where user_id = $1;
`

	result := make([]*models.Withdraw, 0)

	if err := u.conn.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, fmt.Errorf("u.conn.SelectContext:%w", err)
	}

	return result, nil
}
