package mock

import (
	"context"

	"github.com/msksgm/go-simple-api-server/model"
)

type UserService struct {
	CreateUserFn   func(*model.User) error
	AuthenticateFn func() *model.User
}

func (m *UserService) CreateUser(_ context.Context, user *model.User) error {
	return m.CreateUserFn(user)
}

func (m *UserService) Authenticate(_ context.Context, username, password string) (*model.User, error) {
	return m.AuthenticateFn(), nil
}
