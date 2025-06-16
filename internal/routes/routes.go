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
		games.POST("/:code/join", app.PlayersHandler.JoinGameHandler)
		games.GET("/:code", app.GamesHandler.GetGameByCode)
	}

	return router
}
