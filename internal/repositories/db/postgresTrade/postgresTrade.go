package db

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type RepositoryTrade struct {
	client postgresql.Client
	logger *logging.Logger
}

type TradeData struct {
	TradeID        uuid.UUID   `json:"trade_id"`
	UserID         uuid.UUID   `json:"user_id"`
	OfferedItems   []uuid.UUID `json:"offered_items"`
	RequestedItems []uuid.UUID `json:"requested_items"`
}

func NewRepository(logger *logging.Logger) *RepositoryTrade {
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	repo := &RepositoryTrade{
		client: client,
		logger: logger,
	}

	return repo
}

func (r *RepositoryTrade) FindAll(ctx context.Context) ([]TradeData, error) {
	q := `
        SELECT trade_id, user_id, offered_items, requested_items FROM public.trade
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	trades := make([]TradeData, 0)

	for rows.Next() {
		var td TradeData

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.OfferedItems, &td.RequestedItems); err != nil {
			return nil, err
		}

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

// FindOne implements trade.Repository.
func (r *RepositoryTrade) FindOne(ctx context.Context, id string) (TradeData, error) {
	q := `
        SELECT trade_id, user_id, offered_items, requested_items FROM public.trade WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var td TradeData
	err := r.client.QueryRow(ctx, q, id).Scan(&td.TradeID, &td.UserID, &td.OfferedItems, &td.RequestedItems)
	if err != nil {
		return TradeData{}, err
	}

	return td, nil
}
func (r *RepositoryTrade) Create(ctx context.Context, data TradeData) (interface{}, error) {
	q := `
		INSERT INTO public.trade (
			trade_id,
			user_id,
			offered_items,
			requested_items)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3)
		RETURNING trade_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := r.client.QueryRow(ctx, q, data.UserID, data.OfferedItems, data.RequestedItems).Scan(&data.TradeID); err != nil {
		r.logger.Infof("Failed to create trade: %v", data)
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err
	}

	r.logger.Infof("Completed to create trade: %v", data)
	return data.TradeID, nil
}

// FindByItemUUID implements trade.Repository.
func (r *RepositoryTrade) FindByItemUUID(ctx context.Context, itemID string) ([]TradeData, error) {
	q := `
        SELECT trade_id, user_id, offered_items, requested_items FROM public.trade
		WHERE $1 = ANY(offered_items) OR $1 = ANY(requested_items)
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q, itemID)
	if err != nil {
		return nil, err
	}

	trades := make([]TradeData, 0)

	for rows.Next() {
		var td TradeData

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.OfferedItems, &td.RequestedItems); err != nil {
			return nil, err
		}

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

// UpdateByID implements trade.Repository.
func (r *RepositoryTrade) Update(ctx context.Context, tradeID string, offeredItems, requestedItems []uuid.UUID) error {
	q := `
		UPDATE public.trade
		SET 
			offered_items = $1, 
			requested_items = $2
		WHERE 
			trade_id = $3
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, offeredItems, requestedItems, tradeID); err != nil {
		return err
	}

	return nil
}

// FindByUserUUID implements trade.Repository.
func (r *RepositoryTrade) FindByUserUUID(ctx context.Context, userID string) ([]TradeData, error) {
	q := `
        SELECT trade_id, user_id, offered_items, requested_items FROM public.trade
		WHERE user_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	trades := make([]TradeData, 0)

	for rows.Next() {
		var td TradeData

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.OfferedItems, &td.RequestedItems); err != nil {
			return nil, err
		}

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}
func (r *RepositoryTrade) Delete(ctx context.Context, tradeID string) error {
	q := `
		DELETE FROM public.trade
		WHERE trade_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, tradeID); err != nil {
		return err
	}

	return nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
