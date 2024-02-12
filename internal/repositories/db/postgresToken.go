package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
)

type RepositoryToken struct {
	client postgresql.Client
	logger *logging.Logger
}

type TokenData struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func NewRepositoryToken(logger *logging.Logger) *RepositoryToken {
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	return &RepositoryToken{
		client: client,
		logger: logger,
	}
}

func (r *RepositoryToken) GetExpiredTokens(ctx context.Context, now time.Time) ([]TokenData, error) {
	q := `
		SELECT 
			id, 
			user_id, 
			token, 
			expiration_time
		FROM public.user_token
		WHERE 
			expiration_time < $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]TokenData, 0)

	for rows.Next() {
		var token TokenData

		if err := rows.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpirationTime); err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *RepositoryToken) DeleteToken(ctx context.Context, tokenID uuid.UUID) error {
	q := `
		DELETE 
		FROM public.user_token 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, tokenID)

	return err
}
