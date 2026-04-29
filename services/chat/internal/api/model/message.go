package model

import (
	"time"

	"github.com/google/uuid"
)

type SendMessageRequest struct {
	Body string `json:"body" validate:"required,min=1,max=4096"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageListResponse struct {
	Messages []*MessageResponse `json:"messages"`
	Total    int                `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload MessageResponse `json:"payload"`
}
