package model

import "errors"

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrUnAuthorized      = errors.New("unauthorized")
	ErrNotFound          = errors.New("record not found")
)
