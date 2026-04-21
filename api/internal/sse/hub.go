package sse

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Event represents a server-sent event payload.
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client represents a single SSE connection.
type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Ch     chan Event
	Done   chan struct{}
}

// Hub manages SSE client connections and broadcasts events.
type Hub struct {
	mu       sync.RWMutex
	clients  map[*Client]struct{}
	log      zerolog.Logger
	onRemove func(userID uuid.UUID)
}

// NewHub creates a new SSE hub.
func NewHub(log zerolog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		log:     log,
	}
}

// SetOnRemove sets a callback when a client disconnects.
func (h *Hub) SetOnRemove(fn func(userID uuid.UUID)) {
	h.onRemove = fn
}

// Subscribe registers a new SSE client.
func (h *Hub) Subscribe(client *Client) {
	h.mu.Lock()
	h.clients[client] = struct{}{}
	h.mu.Unlock()
	h.log.Debug().Str("user_id", client.UserID.String()).Msg("SSE client subscribed")
}

// Unsubscribe removes a client.
func (h *Hub) Unsubscribe(client *Client) {
	h.mu.Lock()
	delete(h.clients, client)
	count := len(h.clients)
	h.mu.Unlock()

	close(client.Ch)

	if h.onRemove != nil {
		h.onRemove(client.UserID)
	}

	h.log.Debug().Str("user_id", client.UserID.String()).Int("remaining", count).Msg("SSE client unsubscribed")
}

// SendToUser sends an event to all connections of a specific user.
func (h *Hub) SendToUser(userID uuid.UUID, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Ch <- event:
			case <-client.Done:
			default:
				// Channel full, drop event to avoid blocking
			}
		}
	}
}

// Broadcast sends an event to all connected clients.
func (h *Hub) Broadcast(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Ch <- event:
		case <-client.Done:
		default:
		}
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// UserConnectionCount returns the number of connections for a specific user.
func (h *Hub) UserConnectionCount(userID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	count := 0
	for client := range h.clients {
		if client.UserID == userID {
			count++
		}
	}
	return count
}

// MarshalEvent serializes an SSE event to bytes.
func MarshalEvent(event Event) []byte {
	data, _ := json.Marshal(event.Data)
	return []byte("event: " + event.Type + "\ndata: " + string(data) + "\n\n")
}
