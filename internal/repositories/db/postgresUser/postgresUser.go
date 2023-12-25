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
func (r *repository) Create(ctx context.Context, u interface{}) error {
	user := u.(*model.User)
	q := `
		INSERT INTO public.user (
			id, 
			name, 
			email 
		) 
		VALUES (
			gen_random_uuid(), 
			$1, 
			$2
		)
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, user.Username, user.Email).Scan(&user.UserId); err != nil {
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

// Delete implements user.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM public.user WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements user.Repository.
func (r *repository) FindAll(ctx context.Context) (u []interface{}, err error) {
	q := `
        SELECT id, name, email FROM public.user
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)

	for rows.Next() {
		var u model.User

		if err := rows.Scan(&u.UserId, &u.Username, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Преобразование []model.User в []interface{}
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = user
	}

	return result, nil
}

// FindOne implements user.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (interface{}, error) {
	q := `
        SELECT id, name, email FROM public.user WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u model.User
	err := r.client.QueryRow(ctx, q, id).Scan(&u.UserId, &u.Username, &u.Email)
	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

// Update implements user.Repository.
func (r *repository) Update(ctx context.Context, user interface{}) error {
	u := user.(*model.User)
	q := `
		UPDATE public.user 
		SET 
			name = $2, 
			email = $3 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, u.UserId, u.Username, u.Email)
	if err != nil {
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
