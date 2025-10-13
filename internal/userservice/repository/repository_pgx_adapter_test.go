package repository

import (
    "context"
    "testing"
)

func TestPgxDB_Methods_PanicOnNilPool(t *testing.T) {
    ctx := context.Background()
    db := &pgxDB{} // pool is nil intentionally

    tests := []struct{
        name string
        fn   func()
    }{
        {name: "Begin", fn: func() { _, _ = db.Begin(ctx) }},
        {name: "Exec", fn: func() { _, _ = db.Exec(ctx, "SELECT 1") }},
        {name: "Ping", fn: func() { _ = db.Ping(ctx) }},
        {name: "QueryRow", fn: func() { _ = db.QueryRow(ctx, "SELECT 1") }},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r == nil {
                    t.Fatalf("expected panic for %s with nil pool, got none", tc.name)
                }
            }()
            tc.fn()
        })
    }
}

func TestPgxTxWrapper_Methods_PanicOnNilTx(t *testing.T) {
    ctx := context.Background()
    w := &pgxTxWrapper{} // tx is zero value (nil)

    tests := []struct{
        name string
        fn   func()
    }{
        {name: "QueryRow", fn: func() { _ = w.QueryRow(ctx, "SELECT 1") }},
        {name: "Exec", fn: func() { _, _ = w.Exec(ctx, "SELECT 1") }},
        {name: "Commit", fn: func() { _ = w.Commit(ctx) }},
        {name: "Rollback", fn: func() { _ = w.Rollback(ctx) }},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r == nil {
                    t.Fatalf("expected panic for %s with nil tx, got none", tc.name)
                }
            }()
            tc.fn()
        })
    }
}


