package models

type ResponseAccrual struct {
	OrderID string  `json:"order" db:"order_id"`  // Номер заказа
	Status  string  `json:"status" db:"status"`   // Статус расчёта начисления
	Accrual float64 `json:"accrual" db:"accrual"` //  Рассчитанные баллы к начислению, при отсутствии начисления
}
