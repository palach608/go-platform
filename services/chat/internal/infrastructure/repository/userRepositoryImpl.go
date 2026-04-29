package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (pr *PostgresUserRepository) Create(ctx context.Context, user *domainModel.User) (*domainModel.User, error) {
	err := pr.pool.QueryRow(
		ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		user.Username, user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pr *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainModel.User, error) {
	user := &domainModel.User{}
	err := pr.pool.QueryRow(
		ctx,
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domainModel.User, error) {
	user := &domainModel.User{}
	err := pr.pool.QueryRow(
		ctx,
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *PostgresUserRepository) Update(ctx context.Context, user *domainModel.User) (*domainModel.User, error) {
	err := pr.pool.QueryRow(
		ctx,
		"UPDATE users SET username = $1, email = $2, password_hash = $3,  updated_at = NOW() WHERE id = $4 RETURNING username, email, password_hash, updated_at",
		user.Username, user.Email, user.PasswordHash, user.ID,
	).Scan(&user.Username, &user.Email, &user.PasswordHash, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := pr.pool.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
