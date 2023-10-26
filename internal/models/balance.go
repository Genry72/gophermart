package models

import "time"

// Balance ответ на спрос баланса пользователя
type Balance struct {
	Current   float64 `db:"current" json:"current"`     // Количество баллов, доступное пользователю
	Withdrawn float64 `db:"withdrawn" json:"withdrawn"` // Количество списанных баллов (за веь период)
}

// Withdraw запрос на списание средств
type Withdraw struct {
	UserID int64     `db:"user_id" json:"-"`
	Order  string    `db:"order_id" json:"order"`
	Points float64   `db:"points" json:"sum"`
	Date   time.Time `db:"date" json:"processed_at,omitempty"`
}

// UserBalance Баланс пользователя. Таблица user_balance
type UserBalance struct {
	UserID         int64     `db:"user_id"`
	Accruals       float64   `db:"accruals"`        // Сумма начисленных баллов
	Drawal         float64   `db:"drawal"`          // Сумма списанных баллов
	CurrentBalance float64   `db:"current_balance"` // Актуальный баланс
	LastUpdate     time.Time `db:"last_update"`     // Дата последнего обновления
}
