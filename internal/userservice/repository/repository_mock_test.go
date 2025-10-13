package repository

import (
    "context"
    "errors"
    "testing"

    "github.com/jackc/pgx/v5/pgconn"
)

type mockRow struct{ scanErr error }
func (m mockRow) Scan(dest ...any) error { return m.scanErr }

type mockTx struct{
    row Row
    execErr error
    commitErr error
    rolledBack bool
}
func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) Row { return m.row }
func (m *mockTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) { return pgconn.CommandTag{}, m.execErr }
func (m *mockTx) Commit(ctx context.Context) error { return m.commitErr }
func (m *mockTx) Rollback(ctx context.Context) error { m.rolledBack = true; return nil }

type mockDB struct{
    beginTx *mockTx
    beginErr error
}
func (m *mockDB) Begin(ctx context.Context) (Tx, error) { if m.beginErr != nil { return nil, m.beginErr }; return m.beginTx, nil }
func (m *mockDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) { return pgconn.CommandTag{}, nil }
func (m *mockDB) Ping(ctx context.Context) error { return nil }
func (m *mockDB) QueryRow(ctx context.Context, sql string, args ...any) Row { return mockRow{} }

// execErrDB implements DB but always fails Exec
type execErrDB struct{}
func (execErrDB) Begin(ctx context.Context) (Tx, error) { return &mockTx{}, nil }
func (execErrDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) { return pgconn.CommandTag{}, errors.New("exec error") }
func (execErrDB) Ping(ctx context.Context) error { return nil }
func (execErrDB) QueryRow(ctx context.Context, sql string, args ...any) Row { return mockRow{} }

func TestRegister_ConstraintLogin(t *testing.T) {
    pgErr := &pgconn.PgError{ConstraintName: "users_login_idx"}
    tx := &mockTx{row: mockRow{scanErr: pgErr}}
    repo := NewPostgresUserRepository(&mockDB{beginTx: tx})

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    if !errors.Is(err, ErrorLoginTaken) {
        t.Fatalf("expected ErrorLoginTaken, got %v", err)
    }
}

func TestRegister_ConstraintEmail(t *testing.T) {
    pgErr := &pgconn.PgError{ConstraintName: "users_email_idx"}
    tx := &mockTx{row: mockRow{scanErr: pgErr}}
    repo := NewPostgresUserRepository(&mockDB{beginTx: tx})

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    if !errors.Is(err, ErrorEmailTaken) {
        t.Fatalf("expected ErrorEmailTaken, got %v", err)
    }
}

func TestRegister_GenericQueryError(t *testing.T) {
    someErr := errors.New("db fail")
    tx := &mockTx{row: mockRow{scanErr: someErr}}
    repo := NewPostgresUserRepository(&mockDB{beginTx: tx})

    _, err := repo.Register(context.Background(), "login", "email@example.com", "hash", "code")
    if !errors.Is(err, ErrorQueryFailed) {
        t.Fatalf("expected ErrorQueryFailed, got %v", err)
    }
}

func TestWriteVerificationCode_ExecError(t *testing.T) {
    repo := NewPostgresUserRepository(execErrDB{})
    if err := repo.WriteVerificationCode(context.Background(), 1, "123456"); err == nil {
        t.Fatalf("expected Exec error, got nil")
    }
}


