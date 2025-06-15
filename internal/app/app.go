package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Logger *slog.Logger
}

func NewApplication() (*Application, error) {
	loggerHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(loggerHandler)

	app := &Application{
		Logger: logger,
	}

	return app, nil
}

func (app *Application) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"env":     "dev",
		"version": "1.0.0",
	})
}
