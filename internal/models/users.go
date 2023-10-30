package models

import "time"

type CtxKeyUserID struct{}

type User struct {
	UserID       int64     `db:"user_id" json:"userID"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Email        *string   `db:"email" json:"email,omitempty"`
	FirstName    *string   `db:"first_name" json:"firstName,omitempty"`
	LastName     *string   `db:"last_name" json:"lastName,omitempty"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

type UserRegister struct {
	Username  string  `json:"login" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
}
