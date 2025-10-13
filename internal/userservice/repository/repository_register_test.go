package repository

import (
    "context"
    "errors"
    "testing"

    "github.com/jackc/pgx/v5/pgconn"
)

// Covers successful flow: scan userID, exec upsert, commit
func TestRegister_Success(t *testing.T) {
    tx := &mockTx{row: mockRow{scanErr: nil}}
    db := &mockDB{beginTx: tx}
    repo := NewPostgresUserRepository(db)

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    // our mock Exec returns nil error; Commit returns nil by default
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

// Covers commit error path
func TestRegister_CommitError(t *testing.T) {
    tx := &mockTx{row: mockRow{scanErr: nil}, commitErr: errors.New("commit failed")}
    db := &mockDB{beginTx: tx}
    repo := NewPostgresUserRepository(db)

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    if err == nil {
        t.Fatalf("expected commit error, got nil")
    }
}

// Covers exec error when inserting into email_verifications
func TestRegister_UpsertExecError(t *testing.T) {
    tx := &mockTx{row: mockRow{scanErr: nil}, execErr: errors.New("exec failed")}
    db := &mockDB{beginTx: tx}
    repo := NewPostgresUserRepository(db)

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    if err == nil {
        t.Fatalf("expected exec error, got nil")
    }
}

// Ensure constraint mapping still works (regression check)
func TestRegister_ConstraintMapping(t *testing.T) {
    for name, cons := range map[string]error{
        "login": ErrorLoginTaken,
        "email": ErrorEmailTaken,
    } {
        t.Run(name, func(t *testing.T) {
            var cName string
            if name == "login" { cName = "users_login_idx" } else { cName = "users_email_idx" }
            tx := &mockTx{row: mockRow{scanErr: &pgconn.PgError{ConstraintName: cName}}}
            repo := NewPostgresUserRepository(&mockDB{beginTx: tx})
            _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
            if !errors.Is(err, cons) {
                t.Fatalf("expected %v, got %v", cons, err)
            }
        })
    }
}


