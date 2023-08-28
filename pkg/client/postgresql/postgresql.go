package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"mado/pkg/utils"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, username, password, host, port, database string, maxAttempts int) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	// var pool *pgxpool.Pool
	// var err error   //instead of this used named params

	if err := utils.DoWithTrials(func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second); err != nil {
		log.Fatal("error do with postgresql")
	}
	return pool, nil
}
