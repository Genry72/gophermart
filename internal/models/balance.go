package models

type Balance struct {
	Current   float64 `db:"current" json:"current"`     // Количество баллов, доступное пользователю
	Withdrawn float64 `db:"withdrawn" json:"withdrawn"` // Количество списанных баллов (за веь период)
}
