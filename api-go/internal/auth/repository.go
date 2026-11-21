package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailTaken = errors.New("email already registered")

type UserCredentials struct {
	ID           int64
	PasswordHash string
}

type authRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) error
	FindUserByEmail(ctx context.Context, email string) (UserCredentials, error)
	FindPasswordHash(ctx context.Context, userID int64) (string, error)
	UpdatePasswordHash(ctx context.Context, userID int64, currentHash, newHash string) (bool, error)
	RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
}

type PostgresAuthRepo struct {
	db *pgxpool.Pool
}

var _repo authRepository = (*PostgresAuthRepo)(nil)

func NewPostgresAuthRepo(db *pgxpool.Pool) *PostgresAuthRepo {
	return &PostgresAuthRepo{db: db}
}

func (_dbContext *PostgresAuthRepo) CreateUser(
	ctx context.Context, email, passwordHash string,
) error {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
	`

	_, err := _dbContext.db.Exec(ctx, query, email, passwordHash)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) &&
		pgErr.Code == "23505" &&
		pgErr.ConstraintName == "users_email_unique" {
		return ErrEmailTaken
	}

	return err
}

func (_dbContext *PostgresAuthRepo) FindUserByEmail(
	ctx context.Context, email string,
) (UserCredentials, error) {
	const query = `
		SELECT id, password_hash
		FROM users
		WHERE lower(email) = lower($1)
	`

	var user UserCredentials
	err := _dbContext.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.PasswordHash)
	return user, err
}

func (_dbContext *PostgresAuthRepo) FindPasswordHash(
	ctx context.Context, userID int64,
) (string, error) {
	const query = `
		SELECT password_hash
		FROM users
		WHERE id = $1
	`

	var passwordHash string
	err := _dbContext.db.QueryRow(ctx, query, userID).Scan(&passwordHash)
	return passwordHash, err
}

func (_dbContext *PostgresAuthRepo) UpdatePasswordHash(
	ctx context.Context, userID int64, currentHash, newHash string,
) (bool, error) {
	const query = `
		UPDATE users
		SET password_hash = $3
		WHERE id = $1 AND password_hash = $2
	`

	result, err := _dbContext.db.Exec(ctx, query, userID, currentHash, newHash)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() == 1, nil
}

func (_dbContext *PostgresAuthRepo) RevokeToken(
	ctx context.Context, jti string, expiresAt time.Time,
) error {
	const query = `
		INSERT INTO revoked_tokens (jti, expires_at)
		VALUES ($1, $2)
		ON CONFLICT (jti) DO NOTHING
	`

	_, err := _dbContext.db.Exec(ctx, query, jti, expiresAt)
	return err
}

func (_dbContext *PostgresAuthRepo) IsTokenRevoked(
	ctx context.Context, jti string,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM revoked_tokens
			WHERE jti = $1
		)
	`

	var revoked bool
	err := _dbContext.db.QueryRow(ctx, query, jti).Scan(&revoked)
	return revoked, err
}