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

type UserFilter struct {
	ID       *uint
	Username *string

	Limit  int
	Offset int
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashBytes)

	return nil
}

func (u User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

type UserService interface {
	Authenticate(ctx context.Context, username, password string) (*User, error)

	CreateUser(context.Context, *User) error
}
