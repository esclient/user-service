package repository

import (
    "context"
    "errors"
    "os/user"
    "time"

    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

const (
	RegisterUserQuery = `
	INSERT INTO users (login, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	CheckVerificationCodeExistsQuery = `SELECT EXISTS (SELECT 1 FROM email_verifications WHERE user_id = $1)`

	UpsertVerificationCodeQuery = `
	INSERT INTO email_verifications (user_id, code, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id)
	DO UPDATE SET code = EXCLUDED.code, created_at = EXCLUDED.created_at;
	`
)

// Row abstracts the Scan method used from pgx
type Row interface {
    Scan(dest ...any) error
}

// Tx abstracts the subset of pgx.Tx we need
type Tx interface {
    QueryRow(ctx context.Context, sql string, args ...any) Row
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}

// DB abstracts the subset of pgxpool.Pool we need
type DB interface {
    Begin(ctx context.Context) (Tx, error)
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    QueryRow(ctx context.Context, sql string, args ...any) Row
    Ping(ctx context.Context) error
}

type PostgresUserRepository struct {
    db DB
}

func NewDatabaseConnection(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

    applyPoolTunables(config)

    db, err := newPoolWithConfig(ctx, config)
    if err != nil {
        return nil, err
    }

    if err := pingPool(ctx, db); err != nil {
        return nil, err
    }

    return db, nil
}

// applyPoolTunables centralizes connection pool settings for easier unit testing
func applyPoolTunables(config *pgxpool.Config) {
    config.MaxConns = MaxPoolConns
    config.MinConns = MinPoolConns
    config.MaxConnLifetime = MaxConnLifetime
    config.MaxConnIdleTime = MaxConnIdleTime
}

// Indirections for testability (overridden in tests)
var newPoolWithConfig = pgxpool.NewWithConfig
var pingPool = func(ctx context.Context, db *pgxpool.Pool) error { return db.Ping(ctx) }

// NewPostgresUserRepository accepts an abstract DB (useful for tests)
func NewPostgresUserRepository(db DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

// NewPostgresUserRepositoryFromPool wraps a *pgxpool.Pool for production usage
func NewPostgresUserRepositoryFromPool(pool *pgxpool.Pool) *PostgresUserRepository {
    return &PostgresUserRepository{db: &pgxDB{pool: pool}}
}

func (r *PostgresUserRepository) GetByLogin(login string) (*user.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*user.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) WriteVerificationCode(ctx context.Context, userID int64, verificationCode string) error {
	_, err := r.db.Exec(ctx, UpsertVerificationCodeQuery, userID, verificationCode, time.Now())
	return err
}

func (r *PostgresUserRepository) Register(ctx context.Context, login string, email string, hashedPassword string, verificationCode string) (int64, error) {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return -1, err
    }
    defer tx.Rollback(ctx)

    var userID int64
    err = tx.QueryRow(ctx, RegisterUserQuery, login, email, hashedPassword).Scan(&userID)
    if err != nil {
        var pgxErr *pgconn.PgError
        if errors.As(err, &pgxErr) {
            switch pgxErr.ConstraintName {
            case "users_login_idx":
                return -1, ErrorLoginTaken
            case "users_email_idx":
                return -1, ErrorEmailTaken
            }
        }
        return -1, ErrorQueryFailed
    }

    _, err = tx.Exec(ctx, `
        INSERT INTO email_verifications (user_id, code, created_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id)
        DO UPDATE SET code = EXCLUDED.code, created_at = EXCLUDED.created_at
    `, userID, verificationCode, time.Now())
    if err != nil {
        return -1, err
    }

    if err := tx.Commit(ctx); err != nil {
        return -1, err
    }

    return userID, nil
}