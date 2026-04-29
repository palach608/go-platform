package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
}

type RoomResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoomListResponse struct {
	Rooms []*RoomResponse `json:"rooms"`
	Total int             `json:"total"`
}
