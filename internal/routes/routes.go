package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/app"
	"github.com/y3933y3933/ghost-card/internal/ws"
)

func SetRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	router.GET("/api/healthz", app.HealthCheck)

	// games
	games := router.Group("/api/games")
	{
		games.POST("/", app.GamesHandler.CreateGame)
		games.POST("/:code/join", app.PlayersHandler.JoinGameHandler)
		games.GET("/:code", app.GamesHandler.GetGameByCode)
		games.POST("/:code/rounds", app.RoundsHandler.CreateRoundInOrder)
		games.POST("/:code/draws", app.RoundsHandler.DrawCard)
		games.GET("/:code/rounds/current", app.RoundsHandler.GetCurrentRound)

	}

	// ws
	router.GET("/ws", ws.ServeWS)

	return router
}
