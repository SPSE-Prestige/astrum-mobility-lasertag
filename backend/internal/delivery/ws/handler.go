package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	hub             *Hub
	activeBroadcast map[string]bool
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub:             hub,
		activeBroadcast: make(map[string]bool),
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")
	if gameID == "" {
		http.Error(w, "missing game id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := NewClient(conn, gameID, h.hub)
	h.hub.Register(gameID, client)

	// Start event bus broadcaster for this game if not already running
	if !h.activeBroadcast[gameID] {
		h.activeBroadcast[gameID] = true
		h.hub.StartGameBroadcaster(gameID)
	}

	go client.WritePump()
	go client.ReadPump()
}
