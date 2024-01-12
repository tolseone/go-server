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

func NewRepositoryTrade(logger *logging.Logger) *RepositoryTrade {
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

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
