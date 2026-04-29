package hub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 4096

type IncomingMessage struct {
	Body string `json:"body"`
}

type OutgoingMessage struct {
	Type      string `json:"type"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type Client struct {
	Hub      *Hub
	RoomID   string
	UserID   string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(message, &incoming); err != nil || incoming.Body == "" {
			continue
		}

		outgoing := OutgoingMessage{
			Type:     "message",
			RoomID:   c.RoomID,
			UserID:   c.UserID,
			Username: c.Username,
			Body:     incoming.Body,
		}

		payload, _ := json.Marshal(outgoing)

		c.Hub.Broadcast <- &BroadcastMessage{
			RoomID:   c.RoomID,
			UserID:   c.UserID,
			Username: c.Username,
			Body:     incoming.Body,
			Payload:  payload,
			Ctx:      context.Background(),
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("write error: %v", err)
			return
		}
	}
}
