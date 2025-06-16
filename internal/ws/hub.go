// ws/hub.go
package ws

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessageWithRoom
}

type MessageWithRoom struct {
	GameCode string
	Message  WebSocketMessage
}

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		mu:         sync.RWMutex{},
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan MessageWithRoom),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.GameCode] == nil {
				h.rooms[client.GameCode] = make(map[*Client]bool)
			}
			h.rooms[client.GameCode][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.GameCode]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.GameCode]
			h.mu.RUnlock()

			payload, _ := json.Marshal(msg.Message)
			for client := range clients {
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastToGame(code string, msg WebSocketMessage) {
	h.Broadcast <- MessageWithRoom{
		GameCode: code,
		Message:  msg,
	}
}

func (h *Hub) SendToPlayer(code string, targetPlayerID int64, msg WebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.rooms[code]
	payload, _ := json.Marshal(msg)

	for client := range clients {
		if client.PlayerID == targetPlayerID {
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}
