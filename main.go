package main

import (
	"fmt"
	"os"

	"github.com/y3933y3933/ghost-card/internal/app"
	"github.com/y3933y3933/ghost-card/internal/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	r := routes.SetRoutes(app)

	err = r.Run(fmt.Sprintf(":%d", app.Config.Port))
	if err != nil {
		app.Logger.Error("failed to start server", err)
		os.Exit(1)
	}

}
