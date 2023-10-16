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
