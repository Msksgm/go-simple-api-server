package model

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"-"`
	Username     string    `json:"username,omitempty"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Token        string    `json:"token,omitempty"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashBytes)

	return nil
}

type UserService interface {
	CreateUser(context.Context, *User) error
}
