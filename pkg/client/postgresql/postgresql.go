package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttemps int, username, password, host, port, database string) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)

	err = doWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return nil

	}, maxAttemps, 5*time.Second)

	if err != nil {
		log.Fatal("Error do with tries postgresql")
	}

	return pool, nil
}

// for maxAttemps > 0 {
// 	time.Sleep()
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	pool, err := pgxpool.New(ctx, dsn)
// 	if err != nil {
// 		fmt.Print("failed to connect to postgresql")
// 		return
// 		// TODO : add logger
// 	}
// 	return nil
// }

func doWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}