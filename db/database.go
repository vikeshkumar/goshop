package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var poolInitialized bool
var pool *pgxpool.Pool
var log = zap.NewExample()

func Connect(context context.Context, uri string) (*pgxpool.Pool, error) {
	if !poolInitialized {
		p, err := pgxpool.Connect(context, uri)
		p = pool
		poolInitialized = true
		return p, err
	} else {
		return nil, errors.New("connection already established")
	}
}

func Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return pool.Query(ctx, sql, args)
}

func QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return pool.QueryRow(ctx, sql, args)
}

func SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return pool.SendBatch(ctx, b)
}

func Begin(ctx context.Context) (pgx.Tx, error) {
	return pool.Begin(ctx)
}

func BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return pool.BeginTx(ctx, txOptions)
}

func CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
