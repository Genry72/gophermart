package models

type ResponseAccrual struct {
	OrderID string  `json:"order" db:"order_id"`  // Номер заказа
	Status  string  `json:"status" db:"status"`   // Статус расчёта начисления
	Accrual float64 `json:"accrual" db:"accrual"` //  Рассчитанные баллы к начислению, при отсутствии начисления
}

type ResponseAccruals []*ResponseAccrual

// GetOrderIDs получение списка order_id.
func (a ResponseAccruals) GetOrderIDs() []string {
	result := make([]string, 0, len(a))

	for i := range a {
		if a[i] != nil {
			result = append(result, a[i].OrderID)
		}
	}

	return result
}
