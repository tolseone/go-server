package db

import (
	"context"
	"fmt"
	"go-server/internal/models/itemModel"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

// Create implements item.Repository.
func (r *repository) Create(ctx context.Context, item *item.Item) error {
	q := `
		INSERT INTO item (
			item_id, 
			name, 
			rarity, 
			description) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4)
		RETURNING item_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, item.ItemId, item.Name, item.Rarity, item.Description).Scan(&item.ItemId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil
		}
		return err

	}
	return nil

}

// Delete implements item.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements item.Repository.
func (r *repository) FindAll(ctx context.Context) (i []item.Item, err error) {
	q := `
        SELECT item_id, name, rarity, description FROM item
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	items := make([]item.Item, 0)

	for rows.Next() {
		var it item.Item

		if err := rows.Scan(&it.ItemId, &it.Name, &it.Rarity, &it.Description); err != nil {
			return nil, err
		}

		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// FindOne implements item.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (item.Item, error) {
	q := `
        SELECT item_id, name, rarity, description FROM item WHERE item_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var it item.Item
	err := r.client.QueryRow(ctx, q, id).Scan(&it.ItemId, &it.Name, &it.Rarity, &it.Description)
	if err != nil {
		return item.Item{}, err
	}

	return it, nil
}

// Update implements item.Repository.
func (r *repository) Update(ctx context.Context, item item.Item) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) item.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
