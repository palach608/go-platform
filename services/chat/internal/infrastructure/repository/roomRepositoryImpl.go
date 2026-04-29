package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRoomRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRoomRepository(pool *pgxpool.Pool) *PostgresRoomRepository {
	return &PostgresRoomRepository{pool: pool}
}

func (pr *PostgresRoomRepository) Create(ctx context.Context, room *domainModel.Room) (*domainModel.Room, error) {
	err := pr.pool.QueryRow(
		ctx,
		"INSERT INTO rooms (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		room.Name, room.Description, room.OwnerId,
	).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (pr *PostgresRoomRepository) FindAll(ctx context.Context) ([]*domainModel.Room, error) {
	rows, err := pr.pool.Query(ctx,
		"SELECT id, name, description, owner_id, created_at, updated_at FROM rooms")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []*domainModel.Room
	for rows.Next() {
		rm := &domainModel.Room{}
		err := rows.Scan(&rm.ID, &rm.Name, &rm.Description, &rm.OwnerId, &rm.CreatedAt, &rm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, nil
}

func (pr *PostgresRoomRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainModel.Room, error) {
	rm := &domainModel.Room{}
	err := pr.pool.QueryRow(
		ctx,
		"SELECT id, name, description, owner_id, created_at, updated_at FROM rooms WHERE id = $1",
		id,
	).Scan(&rm.ID, &rm.Name, &rm.Description, &rm.OwnerId, &rm.CreatedAt, &rm.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

func (pr *PostgresRoomRepository) Update(ctx context.Context, room *domainModel.Room) (*domainModel.Room, error) {
	err := pr.pool.QueryRow(
		ctx,
		"UPDATE rooms SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 RETURNING name, description, updated_at",
		room.Name, room.Description, room.ID,
	).Scan(&room.Name, &room.Description, &room.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (pr *PostgresRoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := pr.pool.Exec(
		ctx,
		"DELETE FROM rooms WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
