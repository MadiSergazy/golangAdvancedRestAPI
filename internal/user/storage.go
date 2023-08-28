package user

import (
	"context"
)

// if in future we want to change DB
type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) (u []User, err error)
}
