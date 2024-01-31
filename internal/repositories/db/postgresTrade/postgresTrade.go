package db

import (
	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"

	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	Status         string      `json:"status"`
	Date           time.Time   `json:"date"`
	OfferedItems   []TradeItem `json:"offered_items"`
	RequestedItems []TradeItem `json:"requested_items"`
}

type TradeItem struct {
	ItemID     uuid.UUID `json:"item_id"`
	ItemStatus string    `json:"item_status"`
}

func NewRepository(logger *logging.Logger) *RepositoryTrade {
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	return &RepositoryTrade{
		client: client,
		logger: logger,
	}
}

func (r *RepositoryTrade) Create(ctx context.Context, data TradeData) (interface{}, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		_ = tx.Commit(ctx)
	}()

	// 1. Добавить запись в таблицу trade
	qTrade := `
		INSERT INTO public.trade (
			id,
			user_id,
			status,
			date)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			CURRENT_TIMESTAMP)
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(qTrade)))

	if err := tx.QueryRow(ctx, qTrade, data.UserID, data.Status).Scan(&data.TradeID); err != nil {
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

	// 2. Добавить записи в таблицу trade_item
	qTradeItem := `
		INSERT INTO public.trade_item (
			id,
			trade_id,
			item_id,
			item_status)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3)
		RETURNING id
	`

	for _, item := range append(data.OfferedItems, data.RequestedItems...) {
		if _, err := tx.Exec(ctx, qTradeItem, data.TradeID, item.ItemID, item.ItemStatus); err != nil {
			return nil, err
		}
	}

	r.logger.Infof("Completed to create trade: %v", data)
	return data.TradeID, nil
}

func (r *RepositoryTrade) FindAll(ctx context.Context) ([]TradeData, error) {
	q := `
        SELECT 
			id, 
			user_id, 
			status, 
			date 
		FROM public.trade
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	trades := make([]TradeData, 0)

	for rows.Next() {
		var td TradeData

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date); err != nil {
			return nil, err
		}

		offeredItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "offered")
		if err != nil {
			return nil, err
		}

		requestedItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "requested")
		if err != nil {
			return nil, err
		}

		td.OfferedItems = offeredItems
		td.RequestedItems = requestedItems

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) FindOne(ctx context.Context, id string) (TradeData, error) {
	q := `
        SELECT 
			id, 
			user_id, 
			status, 
			date 
		FROM public.trade 
		WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var td TradeData
	err := r.client.QueryRow(ctx, q, id).Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date)
	if err != nil {
		return TradeData{}, err
	}

	offeredItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "offered")
	if err != nil {
		return TradeData{}, err
	}

	requestedItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "requested")
	if err != nil {
		return TradeData{}, err
	}

	td.OfferedItems = offeredItems
	td.RequestedItems = requestedItems

	return td, nil
}

func (r *RepositoryTrade) FindByItemUUID(ctx context.Context, itemID string) ([]TradeData, error) {
	q := `
        SELECT 
			id, 
			user_id, 
			status, 
			date 
		FROM public.trade
		WHERE 
			$1 = ANY(offered_items) 
			OR 
			$1 = ANY(requested_items)
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q, itemID)
	if err != nil {
		return nil, err
	}

	trades := make([]TradeData, 0)

	for rows.Next() {
		var td TradeData

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date); err != nil {
			return nil, err
		}

		offeredItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "offered")
		if err != nil {
			return nil, err
		}

		requestedItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "requested")
		if err != nil {
			return nil, err
		}

		td.OfferedItems = offeredItems
		td.RequestedItems = requestedItems

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) Update(ctx context.Context, tradeID string, offeredItems, requestedItems []uuid.UUID) error {
	q := `
		UPDATE public.trade
		SET 
			offered_items = $1, 
			requested_items = $2
		WHERE 
			id = $3
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, offeredItems, requestedItems, tradeID); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryTrade) FindByUserUUID(ctx context.Context, userID string) ([]TradeData, error) {
	q := `
        SELECT id, user_id, status, date FROM public.trade
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

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date); err != nil {
			return nil, err
		}

		offeredItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "offered")
		if err != nil {
			return nil, err
		}

		requestedItems, err := loadTradeItems(ctx, r.client, r.logger, td.TradeID, "requested")
		if err != nil {
			return nil, err
		}

		td.OfferedItems = offeredItems
		td.RequestedItems = requestedItems

		trades = append(trades, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) Delete(ctx context.Context, tradeID string) error {
	q := `
		DELETE 
		FROM public.trade
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, tradeID); err != nil {
		return err
	}

	return nil
}

func loadTradeItems(ctx context.Context, client postgresql.Client, logger *logging.Logger, tradeID uuid.UUID, itemStatus string) ([]TradeItem, error) {
	q := `
        SELECT
			item_id, 
			item_status 
		FROM public.trade_item 
		WHERE 
			trade_id = $1 
		AND 
			item_status = $2
	`
	logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := client.Query(ctx, q, tradeID, itemStatus)
	if err != nil {
		return nil, err
	}

	var tradeItems []TradeItem

	for rows.Next() {
		var ti TradeItem
		if err := rows.Scan(&ti.ItemID, &ti.ItemStatus); err != nil {
			return nil, err
		}

		tradeItems = append(tradeItems, ti)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tradeItems, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
