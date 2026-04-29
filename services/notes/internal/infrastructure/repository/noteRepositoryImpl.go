package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	domainModel "github.com/palach608/go-platform/services/notes/internal/domain/model"
)

type PostgresNoteRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresNoteRepository(pool *pgxpool.Pool) *PostgresNoteRepository {
	return &PostgresNoteRepository{pool: pool}
}

func (pr *PostgresNoteRepository) Create(note *domainModel.Note) (*domainModel.Note, error) {
	err := pr.pool.QueryRow(
		context.Background(),
		"INSERT INTO notes (user_id, title, body) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		note.UserID, note.Title, note.Body,
	).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (pr *PostgresNoteRepository) Delete(id uuid.UUID) error {
	_, err := pr.pool.Exec(
		context.Background(),
		"DELETE FROM notes WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresNoteRepository) FindAll() ([]*domainModel.Note, error) {
	rows, err := pr.pool.Query(
		context.Background(),
		"SELECT * FROM notes",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []*domainModel.Note
	for rows.Next() {
		note := &domainModel.Note{}
		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (pr *PostgresNoteRepository) FindByID(id uuid.UUID) (*domainModel.Note, error) {
	note := &domainModel.Note{}
	err := pr.pool.QueryRow(
		context.Background(),
		"SELECT * FROM notes WHERE id = $1",
		id,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (pr *PostgresNoteRepository) FindByUserID(userID uuid.UUID) ([]*domainModel.Note, error) {
	rows, err := pr.pool.Query(
		context.Background(),
		"SELECT * FROM notes WHERE user_id = $1", userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []*domainModel.Note
	for rows.Next() {
		note := &domainModel.Note{}
		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (pr *PostgresNoteRepository) Update(note *domainModel.Note) (*domainModel.Note, error) {
	err := pr.pool.QueryRow(
		context.Background(),
		"UPDATE notes SET title = $1, body = $2, updated_at = NOW() WHERE id = $3 RETURNING id, user_id, title, body, created_at, updated_at",
		note.Title, note.Body, note.ID,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}
