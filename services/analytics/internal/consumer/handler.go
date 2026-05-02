package consumer

import (
	"context"
	"encoding/json"
	"log"

	kafkapkg "github.com/palach608/go-platform/pkg/broker"
	chclient "github.com/palach608/go-platform/services/analytics/internal/clickhouse"
)

type EventHandler struct {
	ch *chclient.Client
}

func NewEventHandler(ch *chclient.Client) *EventHandler {
	return &EventHandler{ch: ch}
}

func (h *EventHandler) Handle(ctx context.Context, data []byte) error {
	var event kafkapkg.NoteCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("analytics: inserting event %s for user %s", event.Type, event.Username)

	return h.ch.InsertEvent(ctx, event.Type, event.UserID, event.Username, event.NoteID, event.Title)
}
