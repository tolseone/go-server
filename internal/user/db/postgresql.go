package user

import (
	"context"
	"fmt"
	"go-server/internal/user"
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
func (r *repository) Create(ctx context.Context, user *user.User) error {
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
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
	panic("unimplemented")
}

// FindAll implements user.Repository.
func (r *repository) FindAll(ctx context.Context) (u []user.User, err error) {
	q := `
        SELECT id, name, email FROM public.user
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]user.User, 0)

	for rows.Next() {
		var u user.User

		if err := rows.Scan(&u.UserId, &u.Username, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindOne implements user.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `
        SELECT id, name, email FROM public.user WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&u.UserId, &u.Username, &u.Email)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Update implements user.Repository.
func (r *repository) Update(ctx context.Context, user user.User) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
