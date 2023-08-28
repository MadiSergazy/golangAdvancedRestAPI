package author

import (
	"context"
)

// if in future we want to change DB
type Repository interface {
	Create(ctx context.Context, user *Author) error
	FindOne(ctx context.Context, id string) (Author, error)
	Update(ctx context.Context, user Author) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) (u []Author, err error)
}
