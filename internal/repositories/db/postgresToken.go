package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

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
	UserAgent      string    `json:"user_agent"`
	UserRole       string    `json:"user_role"`
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

func (r *RepositoryToken) Create(ctx context.Context, t interface{}) (interface{}, error) {
	q := `
		INSERT INTO public.user_token ( 
			id, 
			user_id, 
			token, 
			expiration_time,
			user_agent
		)
		VALUES (
			gen_random_uuid(), 
			$1,
			$2,
			$3,
			$4
		) 
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	tokenData := t.(TokenData)

	if err := r.client.QueryRow(ctx, q, tokenData.UserID, tokenData.Token, tokenData.ExpirationTime, tokenData.UserAgent).Scan(&tokenData.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err

	}
	r.logger.Infof("Completed to create token: %v", tokenData)
	return tokenData.ID, nil

}

func (r *RepositoryToken) Update(ctx context.Context, t interface{}) (interface{}, error) {
	q := `
		UPDATE public.user_token 
		SET 
			user_id = $1,
			token = $2,
			expiration_time = $3,
			user_agent = $4
		WHERE 
			id = $5
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	tokenData := t.(TokenData)

	_, err := r.client.Exec(ctx, q, tokenData.UserID, tokenData.Token, tokenData.ExpirationTime, tokenData.UserAgent, tokenData.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err
	}

	r.logger.Infof("Completed to update token: %v", tokenData.ID)
	return tokenData.ID, nil
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

func (r *RepositoryToken) GetTokenByUA(ctx context.Context, ua string) (TokenData, error) {
	q := `
		SELECT 
			id, 
			user_id, 
			token, 
			expiration_time,
			user_agent
		FROM public.user_token
		WHERE 
			user_agent = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var token TokenData
	err := r.client.QueryRow(ctx, q, ua).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpirationTime, &token.UserAgent)
	if err != nil {
		return TokenData{}, err
	}

	return token, nil
}

func (r *RepositoryToken) DeleteToken(ctx context.Context, tokenID uuid.UUID) error {
	q := `
		DELETE 
		FROM public.user_token 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, tokenID); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryToken) DeleteTokenByUserID(ctx context.Context, userID string) error {
	q := `
		DELETE 
		FROM public.user_token 
		WHERE 
			user_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, userID); err != nil {
		return err
	}

	return nil
}
