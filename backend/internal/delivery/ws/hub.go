package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/gorilla/websocket"
)

// Hub manages WebSocket connections per game.
type Hub struct {
	mu       sync.RWMutex
	games    map[string]map[*Client]struct{}
	eventBus domain.EventBus
}

func NewHub(eventBus domain.EventBus) *Hub {
	return &Hub{
		games:    make(map[string]map[*Client]struct{}),
		eventBus: eventBus,
	}
}

func (h *Hub) Register(gameID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.games[gameID] == nil {
		h.games[gameID] = make(map[*Client]struct{})
	}
	h.games[gameID][client] = struct{}{}
}

func (h *Hub) Unregister(gameID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.games[gameID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.games, gameID)
		}
	}
}

func (h *Hub) Broadcast(gameID string, msg domain.WSMessage) {
	h.mu.RLock()
	clients := h.games[gameID]
	h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for c := range clients {
		select {
		case c.send <- data:
		default:
			go func(client *Client) {
				h.Unregister(gameID, client)
				client.conn.Close()
			}(c)
		}
	}
}

// Run subscribes to the event bus and broadcasts messages to connected clients.
func (h *Hub) Run() {
	// This is started per-game when clients connect.
	// The per-game subscription is managed in the Handler.
}

// StartGameBroadcaster subscribes to a game's event bus channel and broadcasts to WS clients.
func (h *Hub) StartGameBroadcaster(gameID string) {
	ch := h.eventBus.Subscribe(gameID)
	go func() {
		for msg := range ch {
			h.Broadcast(gameID, msg)
		}
	}()
}

// Client represents a single WebSocket connection.
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	gameID string
	hub    *Hub
}

func NewClient(conn *websocket.Conn, gameID string, hub *Hub) *Client {
	return &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		gameID: gameID,
		hub:    hub,
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c.gameID, c)
		c.conn.Close()
		close(c.send)
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			return
		}
		// Client messages are not processed; this is a broadcast-only channel.
	}
}
