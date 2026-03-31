package repository

import (
    "context"
    "testing"
)

// Expect panics on nil internals; we just want to exercise adapter code paths without real DB
func TestPgxDB_Methods_PanicOnNilPool(t *testing.T) {
    ctx := context.Background()
    db := &pgxDB{} // nil pool

    tests := []struct{
        name string
        fn   func()
    }{
        {name: "QueryRow", fn: func() { _ = db.QueryRow(ctx, "SELECT 1") }},
        {name: "Exec", fn: func() { _, _ = db.Exec(ctx, "SELECT 1") }},
        {name: "Ping", fn: func() { _ = db.Ping(ctx) }},
        {name: "Begin", fn: func() { _, _ = db.Begin(ctx) }},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r == nil {
                    t.Fatalf("expected panic with nil pool for %s, got none", tc.name)
                }
            }()
            tc.fn()
        })
    }
}

func TestPgxTxWrapper_Methods_PanicOnNilTx(t *testing.T) {
    ctx := context.Background()
    w := &pgxTxWrapper{}

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
                    t.Fatalf("expected panic with nil tx for %s, got none", tc.name)
                }
            }()
            tc.fn()
        })
    }
}


