package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/app"
	"github.com/y3933y3933/joker/internal/ws"
)

func SetRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	router.GET("/api/healthz", app.HealthCheck)

	// games
	games := router.Group("/api/games")
	{
		games.POST("/", app.GamesHandler.CreateGame)
		games.GET("/:code/players", app.PlayersHandler.ListPlayers)
		games.GET("/:code/rounds/current", app.RoundsHandler.GetCurrentRound)
		games.POST("/:code/join", app.PlayersHandler.JoinGame)
		games.POST("/:code/rounds", app.RoundsHandler.CreateRound)
		games.POST("/:code/rounds/:id/draw", app.RoundsHandler.DrawCard)
		games.POST("/:code/rounds/next", app.RoundsHandler.CreateNextRound)
		games.POST("/:code/end", app.RoundsHandler.EndGame)
	}

	// ws
	router.GET("/ws/games/:code", func(c *gin.Context) {
		ws.ServeWS(app.WSHub, c)
	})

	return router
}
