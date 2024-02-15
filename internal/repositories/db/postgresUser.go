package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"go-server/internal/config"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
)

type RepositoryUser struct {
	client postgresql.Client
	logger *logging.Logger
}

type UserData struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role,omitempty"`
}

func NewRepositoryUser(logger *logging.Logger) *RepositoryUser {
	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	return &RepositoryUser{
		client: client,
		logger: logger,
	}

}

func (r *RepositoryUser) Create(ctx context.Context, u interface{}) (interface{}, error) {
	q := `
		INSERT INTO public.user (
			id, 
			name, 
			email,
			password,
			role 
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
	userData := u.(UserData)

	if err := r.client.QueryRow(ctx, q, userData.Username, userData.Email, userData.Password, userData.Role).Scan(&userData.UserId); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err

	}
	r.logger.Infof("Completed to create user: %v", userData)
	return userData.UserId, nil

}

func (r *RepositoryUser) Delete(ctx context.Context, id string) error {
	q := `
		DELETE 
		FROM public.user 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if _, err := r.client.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUser) FindAll(ctx context.Context) ([]UserData, error) {
	q := `
        SELECT 
			id, 
			name, 
			email,
			password,
			role
		FROM public.user
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]UserData, 0)

	for rows.Next() {
		var us UserData

		if err := rows.Scan(&us.UserId, &us.Username, &us.Email, &us.Password, &us.Role); err != nil {
			return nil, err
		}

		users = append(users, us)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RepositoryUser) FindOne(ctx context.Context, id string) (UserData, error) {
	q := `
        SELECT 
			id, 
			name, 
			email,
			password,
			role
		FROM public.user 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u UserData
	err := r.client.QueryRow(ctx, q, id).Scan(&u.UserId, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return UserData{}, err
	}

	return u, nil
}

func (r *RepositoryUser) FindUserByEmail(ctx context.Context, email string) (UserData, error) {
	q := `
        SELECT 
			id, 
			name, 
			email,
			password,
			role 
		FROM public.user 
		WHERE 
			email = $1
    `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u UserData
	err := r.client.QueryRow(ctx, q, email).Scan(&u.UserId, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return UserData{}, err
	}

	return u, nil
}

func (r *RepositoryUser) Update(ctx context.Context, user interface{}) (interface{}, error) {
	q := `
		UPDATE public.user 
		SET 
			name = $2, 
			email = $3,
			password = $4,
			role = $5 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	updatedUser := user.(UserData)

	if _, err := r.client.Exec(ctx, q, updatedUser.UserId, updatedUser.Username, updatedUser.Email, updatedUser.Password, updatedUser.Role); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RepositoryUser) UpdateUserRole(ctx context.Context, id, role string) error {
	q := `
		UPDATE public.user 
		SET 
			role = $2 
		WHERE 
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id, role)
	if err != nil {
		return err
	}

	return nil

}
