package ws

type Room struct {
	Code       string           // 房間代碼
	Clients    map[*Client]bool // 客戶端清單
	Broadcast  chan []byte      // 廣播訊息給所有人
	Register   chan *Client     // 有人加入
	Unregister chan *Client     // 有人離開
}

func NewRoom(code string) *Room {
	return &Room{
		Code:       code,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.Register:
			r.Clients[client] = true
			joinMsg := []byte("🔔 A new player has joined room " + r.Code)
			r.broadcastToAll(joinMsg)

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				leaveMsg := []byte("👋 A player has left room " + r.Code)
				r.broadcastToAll(leaveMsg)
			}

		case message := <-r.Broadcast:
			r.broadcastToAll(message)
		}
	}
}

func (r *Room) broadcastToAll(message []byte) {
	for client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(r.Clients, client)
		}
	}
}
