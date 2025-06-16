package ws

import "sync"

type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

// GetOrCreateRoom 回傳一個已存在的房間，或創建新的
func (h *Hub) GetOrCreateRoom(code string) *Room {
	h.mu.RLock()
	room, exists := h.rooms[code]
	h.mu.RUnlock()

	if exists {
		return room
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// double-check pattern
	if room, exists := h.rooms[code]; exists {
		return room
	}

	room = NewRoom(code)
	h.rooms[code] = room
	go room.Run()
	return room
}
