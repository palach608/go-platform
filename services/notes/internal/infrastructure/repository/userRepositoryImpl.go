package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-notes/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (pr *PostgresUserRepository) Create(user *domainModel.User) (*domainModel.User, error) {
	err := pr.pool.QueryRow(
		context.Background(),
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at",
		user.Username, user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *PostgresUserRepository) Delete(id uuid.UUID) error {
	_, err := pr.pool.Exec(
		context.Background(),
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepository) FindByEmail(email string) (*domainModel.User, error) {
	var user domainModel.User
	err := pr.pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE email = $1", email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pr *PostgresUserRepository) FindByID(id uuid.UUID) (*domainModel.User, error) {
	var user domainModel.User
	err := pr.pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
