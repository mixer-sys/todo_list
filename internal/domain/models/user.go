package models

import (
	"time"
)

type User struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Password   string    `json:"-"`
	Enable2FA  bool      `json:"enable_2fa"`
	TGUserName string    `json:"tg_username"`
}

type UserWithCode struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Password   string    `json:"-"`
	Enable2FA  bool      `json:"enable_2fa"`
	TGUserName string    `json:"tg_username"`
	Code       string    `json:"code"`
}
