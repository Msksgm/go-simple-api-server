package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/msksgm/go-simple-api-server/model"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (us *UserService) CreateUser(ctx context.Context, user *model.User) error {
	tx, err := us.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}
	return tx.Commit()
}

func createUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2) RETURNING id, created_at, updated_at
	`
	args := []interface{}{user.Username, user.PasswordHash}
	err := tx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return model.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}
