package hub

import "github.com/google/uuid"

type Room struct {
	id      string
	clients map[*Client]bool
}

func NewRoom(id string) *Room {
	return &Room{
		id:      id,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Clients() map[*Client]bool {
	return r.clients
}

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
