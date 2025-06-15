package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/app"
)

func SetRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", app.HealthCheck)

	return router
}
