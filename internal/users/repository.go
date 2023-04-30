package users

import (
	"context"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)
	AddUser(ctx context.Context, user User) (int64, error)
	UpdateByID(ctx context.Context, id int64, user User) error
	DeleteUser(ctx context.Context, id int64)
}
