package clickhouse

import (
	"context"
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	db *sql.DB
}

func NewClient(dsn string) (*Client, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	c := &Client{db: db}
	if err := c.migrate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) migrate() error {
	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			event_type  String,
			user_id     String,
			username    String,
			entity_id   String,
			title       String,
			created_at  DateTime DEFAULT now()
		) ENGINE = MergeTree()
		ORDER BY (event_type, created_at)
	`)
	return err
}

func (c *Client) InsertEvent(ctx context.Context, eventType, userID, username, entityID, title string) error {
	_, err := c.db.ExecContext(ctx,
		`INSERT INTO events (event_type, user_id, username, entity_id, title) VALUES (?, ?, ?, ?, ?)`,
		eventType, userID, username, entityID, title,
	)
	return err
}

func (c *Client) TopUsers(ctx context.Context) ([]TopUser, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT username, count() as total
		FROM events
		WHERE created_at >= now() - INTERVAL 7 DAY
		GROUP BY username
		ORDER BY total DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TopUser
	for rows.Next() {
		var u TopUser
		if err := rows.Scan(&u.Username, &u.Total); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (c *Client) NotesPerDay(ctx context.Context) ([]NotesDay, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT toDate(created_at) as day, count() as total
		FROM events
		WHERE event_type = 'note.created'
		  AND created_at >= now() - INTERVAL 7 DAY
		GROUP BY day
		ORDER BY day
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []NotesDay
	for rows.Next() {
		var d NotesDay
		if err := rows.Scan(&d.Day, &d.Total); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

type TopUser struct {
	Username string `json:"username"`
	Total    uint64 `json:"total"`
}

type NotesDay struct {
	Day   string `json:"day"`
	Total uint64 `json:"total"`
}
