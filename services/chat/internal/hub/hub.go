package hub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SilverName608/go-chat/internal/domain/service"
)

type BroadcastMessage struct {
	RoomID   string
	UserID   string
	Username string
	Body     string
	Payload  []byte
	Ctx      context.Context
}

type SystemMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Body     string `json:"body"`
}

type Hub struct {
	Rooms          map[string]*Room
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan *BroadcastMessage
	MessageService service.MessageService
}

func NewHub(msgService service.MessageService) *Hub {
	return &Hub{
		Rooms:          make(map[string]*Room),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan *BroadcastMessage),
		MessageService: msgService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room, ok := h.Rooms[client.RoomID]
			if !ok {
				room = NewRoom(client.RoomID)
				h.Rooms[client.RoomID] = room
			}
			room.clients[client] = true
			log.Printf("client %s joined room %s", client.Username, client.RoomID)
			h.sendSystem(room, client, "join", client.Username)

		case client := <-h.Unregister:
			room, ok := h.Rooms[client.RoomID]
			if !ok {
				continue
			}
			if _, exists := room.clients[client]; exists {
				delete(room.clients, client)
				close(client.Send)
				log.Printf("client %s left room %s", client.Username, client.RoomID)
				h.sendSystem(room, nil, "leave", client.Username)
				if len(room.clients) == 0 {
					delete(h.Rooms, client.RoomID)
				}
			}

		case msg := <-h.Broadcast:
			room, ok := h.Rooms[msg.RoomID]
			if !ok {
				continue
			}
			if h.MessageService != nil && msg.Body != "" {
				roomUUID := mustParseUUID(msg.RoomID)
				userUUID := mustParseUUID(msg.UserID)
				_, err := h.MessageService.Create(msg.Ctx, roomUUID, userUUID, msg.Body)
				if err != nil {
					log.Printf("failed to save message: %v", err)
				}
			}
			for client := range room.clients {
				select {
				case client.Send <- msg.Payload:
				default:
					close(client.Send)
					delete(room.clients, client)
				}
			}
		}
	}
}

func (h *Hub) sendSystem(room *Room, exclude *Client, msgType, username string) {
	msg := SystemMessage{
		Type:     msgType,
		Username: username,
		Body:     username,
	}
	payload, _ := json.Marshal(msg)
	for client := range room.clients {
		if client == exclude {
			continue
		}
		select {
		case client.Send <- payload:
		default:
		}
	}
}
