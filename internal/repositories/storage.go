package storage

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, entity interface{}) error
	FindAll(ctx context.Context) ([]interface{}, error)
	FindOne(ctx context.Context, id string) (interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id string) error
}
