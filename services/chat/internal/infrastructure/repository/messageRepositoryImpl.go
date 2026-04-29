package repository

import (
	"context"

	domainModel "github.com/SilverName608/go-chat/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMessageRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMessageRepository(pool *pgxpool.Pool) *PostgresMessageRepository {
	return &PostgresMessageRepository{pool: pool}
}

func (pr *PostgresMessageRepository) Create(ctx context.Context, message *domainModel.Message) (*domainModel.Message, error) {
	err := pr.pool.QueryRow(
		ctx,
		`INSERT INTO messages (room_id, user_id, body) VALUES ($1, $2, $3) RETURNING id, created_at`,
		message.RoomId, message.UserId, message.Body,
	).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (pr *PostgresMessageRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*domainModel.Message, error) {
	query := `
		SELECT m.id, m.room_id, m.user_id, u.username, m.body, m.created_at
		FROM messages m
		JOIN users u ON u.id = m.user_id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := pr.pool.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domainModel.Message
	for rows.Next() {
		msg := &domainModel.Message{}
		err := rows.Scan(&msg.ID, &msg.RoomId, &msg.UserId, &msg.Username, &msg.Body, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (pr *PostgresMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := pr.pool.Exec(
		ctx,
		"DELETE FROM messages WHERE id = $1",
		id,
	)
	return err
}
