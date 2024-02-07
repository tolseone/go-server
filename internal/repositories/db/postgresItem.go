package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"

)

type RepositoryItem struct {
	client postgresql.Client
	logger *logging.Logger
}

type ItemData struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name"`
	Rarity      string    `json:"rarity"`
	Description string    `json:"description,omitempty"`
}

func NewRepositoryItem(logger *logging.Logger) *RepositoryItem {
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	return &RepositoryItem{
		client: client,
		logger: logger,
	}
}

func (r *RepositoryItem) Create(ctx context.Context, i interface{}) (interface{}, error) {
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
	itemData := i.(ItemData)

	if err := r.client.QueryRow(ctx, q, itemData.Name, itemData.Rarity, itemData.Description).Scan(&itemData.ItemId); err != nil {
		r.logger.Infof("Failed to create item: %v", itemData)
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err
	}
	r.logger.Infof("Completed to create item: %v", itemData)
	return itemData.ItemId, nil

}

func (r *RepositoryItem) Delete(ctx context.Context, id string) error {
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

func (r *RepositoryItem) FindAll(ctx context.Context) ([]ItemData, error) {
	q := `
        SELECT 
			id, 
			name, 
			rarity, 
			description 
		FROM public.item
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	items := make([]ItemData, 0)

	for rows.Next() {
		var it ItemData

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

func (r *RepositoryItem) FindOne(ctx context.Context, id string) (ItemData, error) {
	q := `
        SELECT 
			id, 
			name, 
			rarity, 
			description 
		FROM public.item 
		WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var it ItemData
	err := r.client.QueryRow(ctx, q, id).Scan(&it.ItemId, &it.Name, &it.Rarity, &it.Description)
	if err != nil {
		return ItemData{}, err
	}

	return it, nil
}

func (r *RepositoryItem) Update(ctx context.Context, item interface{}) (interface{}, error) {
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

	updatedItem := item.(ItemData)

	if _, err := r.client.Exec(ctx, q, updatedItem.Name, updatedItem.Rarity, updatedItem.Description, updatedItem.ItemId); err != nil {
		return nil, err
	}

	return nil, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
