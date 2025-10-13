package repository

import (
    "context"
    "errors"
    "testing"

    "github.com/jackc/pgx/v5/pgxpool"
)

// Covers ParseConfig error branch (pure function of input)
func TestNewDatabaseConnection_ParseConfigError(t *testing.T) {
    ctx := context.Background()
    if db, err := NewDatabaseConnection(ctx, "not-a-valid-url"); err == nil {
        if db != nil { db.Close() }
        t.Fatalf("expected parse config error, got nil")
    }
}

// Covers applyPoolTunables, newPoolWithConfig and pingPool via stubs (no real DB)
func TestNewDatabaseConnection_UsesTunables_And_Ping(t *testing.T) {
    ctx := context.Background()

    // Stub constructors
    calledNew := false
    calledPing := false

    // Fake pool implementing only Close and Ping through our pingPool stub
    fakePool := &pgxpool.Pool{}

    origNew := newPoolWithConfig
    origPing := pingPool
    t.Cleanup(func(){ newPoolWithConfig = origNew; pingPool = origPing })

    newPoolWithConfig = func(ctx context.Context, cfg *pgxpool.Config) (*pgxpool.Pool, error) {
        calledNew = true
        // Assert tunables were applied
        if cfg.MaxConns != MaxPoolConns || cfg.MinConns != MinPoolConns || cfg.MaxConnLifetime != MaxConnLifetime || cfg.MaxConnIdleTime != MaxConnIdleTime {
            return nil, errors.New("tunables not applied")
        }
        return fakePool, nil
    }
    pingPool = func(ctx context.Context, db *pgxpool.Pool) error {
        calledPing = true
        return nil
    }

    dsn := "postgres://user:pass@host:5432/dbname?sslmode=disable"
    db, err := NewDatabaseConnection(ctx, dsn)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if db != fakePool { t.Fatalf("unexpected pool instance returned") }
    if !calledNew { t.Fatalf("newPoolWithConfig was not called") }
    if !calledPing { t.Fatalf("pingPool was not called") }
}


