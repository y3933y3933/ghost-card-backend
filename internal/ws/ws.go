package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/y3933y3933/ghost-card/internal/api"
)

// 建立一個 upgrader：把 HTTP 升級成 WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許所有來源（開發階段用）
	},
}

// ServeWS 是 WebSocket 的入口點
func ServeWS(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		api.BadRequest(c, "missing game code")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.InternalServerError(c, "failed to upgrade to websocket")
		return
	}

	msg := map[string]string{
		"message": "WebSocket connected to game " + code,
	}

	conn.WriteJSON(msg)
}
