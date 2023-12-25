package db

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/models"
	"go-server/internal/repositories"
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

// Create implements user.Repository.
// Create implements item.Repository.
func (r *repository) Create(ctx context.Context, i interface{}) error {
	item := i.(*model.Item)
	q := `
		INSERT INTO public.item (
			id, 
			name, 
			rarity, 
			description) 
		VALUES (
			gen_random_uuid(), 
			$1, 
			$2, 
			$3)
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, item.Name, item.Rarity, item.Description).Scan(&item.ItemId); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil

}

// Delete implements item.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM public.item
		WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}

// FindAll implements item.Repository.
func (r *repository) FindAll(ctx context.Context) (i []interface{}, err error) {
	q := `
        SELECT id, name, rarity, description FROM public.item
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	items := make([]model.Item, 0)

	for rows.Next() {
		var it model.Item

		if err := rows.Scan(&it.ItemId, &it.Name, &it.Rarity, &it.Description); err != nil {
			return nil, err
		}

		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Преобразование []model.Item в []interface{}
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item
	}

	return result, nil
}

// FindOne implements item.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (interface{}, error) {
	q := `
        SELECT id, name, rarity, description FROM public.item WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var it model.Item
	err := r.client.QueryRow(ctx, q, id).Scan(&it.ItemId, &it.Name, &it.Rarity, &it.Description)
	if err != nil {
		return model.Item{}, err
	}

	return it, nil
}

// Update implements item.Repository.
func (r *repository) Update(ctx context.Context, item interface{}) error {
	updatedItem := item.(*model.Item)
	q := `
		UPDATE public.item
		SET 
			name = $1, 
			rarity = $2, 
			description = $3
		WHERE 
			id = $4
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, updatedItem.Name, updatedItem.Rarity, updatedItem.Description, updatedItem.ItemId); err != nil {
		return err
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) storage.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
