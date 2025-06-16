package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/app"
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

	}

	return router
}
