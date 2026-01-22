package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ChatMessage struct {
	Type string `json:"type"` // message/join/leave
	Room string `json:"room"`
	From string `json:"from"`
	Text string `json:"text"`
	TS   int64  `json:"ts"`
}

type Client struct {
	id   string
	room string
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]struct{}),
	}
}

func (h *Hub) Join(room string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[room]; !ok {
		h.rooms[room] = make(map[*Client]struct{})
	}
	h.rooms[room][c] = struct{}{}
}

func (h *Hub) Leave(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomClients, ok := h.rooms[c.room]
	if !ok {
		return
	}
	delete(roomClients, c)
	close(c.send)

	// cleanup empty room
	if len(roomClients) == 0 {
		delete(h.rooms, c.room)
	}
}

func (h *Hub) Broadcast(room string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	roomClients, ok := h.rooms[room]
	if !ok {
		return
	}
	for c := range roomClients {
		// non-blocking send, drop if client is slow
		select {
		case c.send <- payload:
		default:
			// slow client -> you can choose to kick them
			// here we just drop the message
		}
	}
}

func main() {
	app := fiber.New()
	hub := NewHub()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		// ---- 1) Get params (room/user) ----
		room := conn.Query("room", "general")
		user := conn.Query("user", "anonymous")

		client := &Client{
			id:   user,
			room: room,
			conn: conn,
			send: make(chan []byte, 32),
		}

		// ---- 2) Join room ----
		hub.Join(room, client)

		// announce join
		joinMsg := ChatMessage{
			Type: "join",
			Room: room,
			From: user,
			Text: user + " joined",
			TS:   time.Now().UnixMilli(),
		}
		b, _ := json.Marshal(joinMsg)
		hub.Broadcast(room, b)

		// ---- 3) Start writer goroutine (single writer per conn) ----
		done := make(chan struct{})
		go func() {
			defer close(done)
			for msg := range client.send {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
		}()

		// ---- 4) Read loop (this goroutine is the reader) ----
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Expect JSON from client
			var incoming ChatMessage
			if err := json.Unmarshal(raw, &incoming); err != nil {
				// If not JSON, treat as plain text
				incoming = ChatMessage{
					Type: "message",
					Room: room,
					From: user,
					Text: string(raw),
					TS:   time.Now().UnixMilli(),
				}
			}

			// Force server truth
			incoming.Type = "message"
			incoming.Room = room
			incoming.From = user
			incoming.TS = time.Now().UnixMilli()

			out, _ := json.Marshal(incoming)
			hub.Broadcast(room, out)
		}

		// ---- 5) Cleanup ----
		hub.Leave(client)

		leaveMsg := ChatMessage{
			Type: "leave",
			Room: room,
			From: user,
			Text: user + " left",
			TS:   time.Now().UnixMilli(),
		}
		lb, _ := json.Marshal(leaveMsg)
		hub.Broadcast(room, lb)

		// wait writer to finish
		<-done
	}))

	log.Fatal(app.Listen(":8080"))
}
