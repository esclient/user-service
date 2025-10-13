package repository

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

type pgxDB struct {
    pool *pgxpool.Pool
}

func (p *pgxDB) Begin(ctx context.Context) (Tx, error) {
    tx, err := p.pool.Begin(ctx)
    if err != nil {
        return nil, err
    }
    return &pgxTxWrapper{tx: tx}, nil
}

func (p *pgxDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
    return p.pool.Exec(ctx, sql, args...)
}

func (p *pgxDB) Ping(ctx context.Context) error {
    return p.pool.Ping(ctx)
}

func (p *pgxDB) QueryRow(ctx context.Context, sql string, args ...any) Row {
    return p.pool.QueryRow(ctx, sql, args...)
}

// pgxpool.Pool.Begin returns pgx.Tx
type pgxTxWrapper struct { tx pgx.Tx }

func (t *pgxTxWrapper) QueryRow(ctx context.Context, sql string, args ...any) Row { return t.tx.QueryRow(ctx, sql, args...) }
func (t *pgxTxWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) { return t.tx.Exec(ctx, sql, args...) }
func (t *pgxTxWrapper) Commit(ctx context.Context) error { return t.tx.Commit(ctx) }
func (t *pgxTxWrapper) Rollback(ctx context.Context) error { return t.tx.Rollback(ctx) }


