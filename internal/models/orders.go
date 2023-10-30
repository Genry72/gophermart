package models

import (
	"time"
)

const (
	OrderStatusNew        = "NEW"        // Заказ зарегистрирован, но вознаграждение не рассчитано
	OrderStatusRegistered = "REGISTERED" // Заказ зарегистрирован, но вознаграждение не рассчитано
	OrderStatusInvalid    = "INVALID"    // Заказ не принят к расчёту, и вознаграждение не будет начислено
	OrderStatusProcessing = "PROCESSING" // Расчёт начисления в процессе
	OrderStatusProcessed  = "PROCESSED"  // Расчёт начисления окончен
)

// Order структура таблицы orders
type Order struct {
	OrderID   string    `db:"order_id" json:"number"`
	UserID    int64     `db:"user_id" json:"-"`
	Status    string    `db:"status" json:"status"`
	Accrual   float64   `db:"accrual" json:"accrual"`
	CreatedAt time.Time `db:"created_at" json:"uploaded_at"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
