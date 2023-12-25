package item

import "context"

type Repository interface {
	Create(ctx context.Context, item *Item) error
	FindAll(ctx context.Context) (i []Item, err error)
	FindOne(ctx context.Context, id string) (Item, error)
	Update(ctx context.Context, item Item) error
	Delete(ctx context.Context, id string) error
}
