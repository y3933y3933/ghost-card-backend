package ws

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 你可以加 origin 檢查
	},
}

func ServeWS(hub *Hub, c *gin.Context) {
	gameCode := c.Param("code")

	playerIDStr := c.Query("player_id")
	playerID, _ := strconv.ParseInt(playerIDStr, 10, 64)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		GameCode: gameCode,
		Hub:      hub,
		PlayerID: playerID,
	}

	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
