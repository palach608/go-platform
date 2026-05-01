package kafkahandler

import (
	"context"
	"encoding/json"
	"log"

	kafkapkg "github.com/palach608/go-platform/pkg/broker"
	"github.com/palach608/go-platform/services/chat/internal/hub"
)

type NoteEventHandler struct {
	hub *hub.Hub
}

func NewNoteEventHandler(h *hub.Hub) *NoteEventHandler {
	return &NoteEventHandler{hub: h}
}

func (h *NoteEventHandler) Handle(ctx context.Context, data []byte) error {
	var event kafkapkg.NoteCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("received note.created event: user=%s title=%s", event.Username, event.Title)

	msg := hub.SystemMessage{
		Type:     "note.created",
		Username: event.Username,
		Body:     event.Username + " создал заметку: " + event.Title,
	}

	payload, _ := json.Marshal(msg)

	for _, room := range h.hub.Rooms {
		for client := range room.Clients() {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}

	return nil
}
