package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 記得未來補上 Origin 驗證
		return true
	},
}

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		playerID := c.Param("player_id")
		if code == "" || playerID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing game code or player ID"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		room := hub.GetOrCreateRoom(code)

		client := &Client{
			ID:   playerID,
			Conn: conn,
			Send: make(chan []byte, 256),
			Room: room,
		}

		// 將 client 加入房間
		room.Register <- client

		// 啟動 goroutine 負責處理這個 client 的發送與接收
		go client.WritePump()
		go client.ReadPump()

		log.Println("✅ WebSocket connected: room =", code, "player =", playerID)
	}
}
