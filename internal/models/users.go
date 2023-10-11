package models

import "time"

type CtxKeyUserID struct{}

type User struct {
	UserID       int64     `db:"user_id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Email        *string   `db:"email"`
	FirstName    *string   `db:"first_name"`
	LastName     *string   `db:"last_name"`
	Phone        *string   `db:"phone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserRegister struct {
	Username  string  `json:"login" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
}
