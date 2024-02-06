package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"

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

	tradeID, err := r.createTrade(ctx, tx, data)
	if err != nil {
		r.logger.Infof("Failed to create trade: %v", data)
		return nil, err
	}

	if err := r.createTradeItems(ctx, tx, tradeID, append(data.OfferedItems, data.RequestedItems...)); err != nil {
		return nil, err
	}

	r.logger.Infof("Completed to create trade: %v", data)
	return tradeID, nil
}

func (r *RepositoryTrade) FindAll(ctx context.Context) ([]TradeData, error) {
	q := `
        SELECT 
			t.id,
			t.user_id,
			t.status,
			t.date,
			ti.item_id,
			ti.item_status
		FROM public.trade t
		LEFT JOIN public.trade_item ti ON t.id = ti.trade_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tradesMap := make(map[uuid.UUID]TradeData)
	for rows.Next() {
		var td TradeData
		var itemID uuid.UUID
		var itemStatus string

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date, &itemID, &itemStatus); err != nil {
			return nil, err
		}

		if itemID != uuid.Nil {
			item := TradeItem{ItemID: itemID, ItemStatus: itemStatus}
			if existingTrade, ok := tradesMap[td.TradeID]; ok {
				if item.ItemStatus == "offered" {
					existingTrade.OfferedItems = append(existingTrade.OfferedItems, item)
				} else if item.ItemStatus == "requested" {
					existingTrade.RequestedItems = append(existingTrade.RequestedItems, item)
				} else {
					r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
				}
				tradesMap[td.TradeID] = existingTrade

			} else {
				if item.ItemStatus == "offered" {
					td.OfferedItems = []TradeItem{item}
				} else if item.ItemStatus == "requested" {
					td.RequestedItems = []TradeItem{item}
				} else {
					r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
				}
				tradesMap[td.TradeID] = td
			}
		}
	}

	var trades []TradeData
	for _, trade := range tradesMap {
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) FindOne(ctx context.Context, tradeID string) (TradeData, error) {
	q := `
        SELECT 
			t.id,
			t.user_id,
			t.status,
			t.date,
			ti.item_id,
			ti.item_status
		FROM public.trade t
		JOIN public.trade_item ti 
		ON 
			t.id = ti.trade_id
		WHERE
			t.id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q, tradeID)
	if err != nil {
		return TradeData{}, err
	}

	defer rows.Close()

	var trade TradeData
	for rows.Next() {
		var itemID uuid.UUID
		var itemStatus string

		if err := rows.Scan(&trade.TradeID, &trade.UserID, &trade.Status, &trade.Date, &itemID, &itemStatus); err != nil {
			return TradeData{}, err
		}

		item := TradeItem{ItemID: itemID, ItemStatus: itemStatus}
		if item.ItemStatus == "offered" {
			trade.OfferedItems = append(trade.OfferedItems, item)
		} else if item.ItemStatus == "requested" {
			trade.RequestedItems = append(trade.RequestedItems, item)
		} else {
			r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
		}
	}

	if err := rows.Err(); err != nil {
		return TradeData{}, err
	}

	return trade, nil
}

func (r *RepositoryTrade) FindByItemUUID(ctx context.Context, itemID string) ([]TradeData, error) {
	q := `
		SELECT 
			t.id,
    		t.user_id,
    		t.status,
			t.date,
			ti.item_id,
			ti.item_status
		FROM public.trade t 
		JOIN public.trade_item ti ON t.id = ti.trade_id
		WHERE EXISTS (
			SELECT 
				1
			FROM public.trade_item ti_sub
			WHERE 
				ti_sub.trade_id = t.id
			AND 
				ti_sub.item_id = $1
			)
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q, itemID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tradesMap := make(map[uuid.UUID]TradeData)
	for rows.Next() {
		var td TradeData
		var itemID uuid.UUID
		var itemStatus string

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date, &itemID, &itemStatus); err != nil {
			return nil, err
		}

		item := TradeItem{ItemID: itemID, ItemStatus: itemStatus}
		if existingTrade, ok := tradesMap[td.TradeID]; ok {
			if item.ItemStatus == "offered" {
				existingTrade.OfferedItems = append(existingTrade.OfferedItems, item)
			} else if item.ItemStatus == "requested" {
				existingTrade.RequestedItems = append(existingTrade.RequestedItems, item)
			} else {
				r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
			}
			tradesMap[td.TradeID] = existingTrade

		} else {
			if item.ItemStatus == "offered" {
				td.OfferedItems = []TradeItem{item}
			} else if item.ItemStatus == "requested" {
				td.RequestedItems = []TradeItem{item}
			} else {
				r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
			}
			tradesMap[td.TradeID] = td
		}
	}

	var trades []TradeData
	for _, trade := range tradesMap {
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) Update(ctx context.Context, trade interface{}) (interface{}, error) {
	updatedTrade := trade.(TradeData)

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

	if err := r.updateTrade(ctx, tx, updatedTrade); err != nil {
		r.logger.Infof("Failed to update trade: %v", updatedTrade)
		return nil, err
	}

	if err := r.updateTradeItems(ctx, tx, updatedTrade.TradeID, append(updatedTrade.OfferedItems, updatedTrade.RequestedItems...)); err != nil {
		return nil, err
	}

	r.logger.Infof("Completed to update trade: %v", updatedTrade)
	return nil, nil
}

func (r *RepositoryTrade) GetTradesByUserUUID(ctx context.Context, userID string) ([]TradeData, error) {
	q := `
        SELECT 
			t.id,
			t.user_id,
			t.status,
			t.date,
			ti.item_id,
			ti.item_status
		FROM public.trade t 
		JOIN public.trade_item ti 
		ON 
			t.id = ti.trade_id
		WHERE 
			t.user_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tradesMap := make(map[uuid.UUID]TradeData)
	for rows.Next() {
		var td TradeData
		var itemID uuid.UUID
		var itemStatus string

		if err := rows.Scan(&td.TradeID, &td.UserID, &td.Status, &td.Date, &itemID, &itemStatus); err != nil {
			return nil, err
		}

		item := TradeItem{ItemID: itemID, ItemStatus: itemStatus}
		if existingTrade, ok := tradesMap[td.TradeID]; ok {
			if item.ItemStatus == "offered" {
				existingTrade.OfferedItems = append(existingTrade.OfferedItems, item)
			} else if item.ItemStatus == "requested" {
				existingTrade.RequestedItems = append(existingTrade.RequestedItems, item)
			} else {
				r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
			}
			tradesMap[td.TradeID] = existingTrade

		} else {
			if item.ItemStatus == "offered" {
				td.OfferedItems = []TradeItem{item}
			} else if item.ItemStatus == "requested" {
				td.RequestedItems = []TradeItem{item}
			} else {
				r.logger.Fatalf("Item status %s is not supported", item.ItemStatus)
			}
			tradesMap[td.TradeID] = td
		}
	}

	var trades []TradeData
	for _, trade := range tradesMap {
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

func (r *RepositoryTrade) createTrade(ctx context.Context, tx pgx.Tx, data TradeData) (uuid.UUID, error) {
	q := `
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
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := tx.QueryRow(ctx, q, data.UserID, data.Status).Scan(&data.TradeID); err != nil {
		return uuid.Nil, err
	}

	return data.TradeID, nil
}

func (r *RepositoryTrade) createTradeItems(ctx context.Context, tx pgx.Tx, tradeID uuid.UUID, items []TradeItem) error {
	q := `
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

	for _, item := range items {
		if _, err := tx.Exec(ctx, q, tradeID, item.ItemID, item.ItemStatus); err != nil {
			return err
		}
	}

	return nil
}

func (r *RepositoryTrade) updateTrade(ctx context.Context, tx pgx.Tx, data TradeData) error {
	q := `
		UPDATE public.trade
		SET
			user_id = $1,
			status = $2,
			date = $3
		WHERE
			id = $4
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := tx.Exec(ctx, q, data.UserID, data.Status, data.Date, data.TradeID); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryTrade) updateTradeItems(ctx context.Context, tx pgx.Tx, tradeID uuid.UUID, items []TradeItem) error {
	q := `
		DELETE FROM public.trade_item 
		WHERE 
			trade_id = $1
	`
	if _, err := tx.Exec(ctx, q, tradeID); err != nil {
		return err
	}

	q = `
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

	for _, item := range items {
		if _, err := tx.Exec(ctx, q, tradeID, item.ItemID, item.ItemStatus); err != nil {
			return err
		}
	}

	return nil
}

func (r *RepositoryTrade) Delete(ctx context.Context, tradeID string) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		_ = tx.Commit(ctx)
	}()

	if err := r.deleteTradeItems(ctx, tx, tradeID); err != nil {
		return err
	}

	if err := r.deleteTrade(ctx, tx, tradeID); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryTrade) deleteTrade(ctx context.Context, tx pgx.Tx, tradeID string) error {
	q := `
        DELETE FROM public.trade
        WHERE id = $1
    `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := tx.Exec(ctx, q, tradeID); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryTrade) deleteTradeItems(ctx context.Context, tx pgx.Tx, tradeID string) error {
	q := `
        DELETE FROM public.trade_item
        WHERE trade_id = $1
    `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := tx.Exec(ctx, q, tradeID); err != nil {
		return err
	}

	return nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
