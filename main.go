package main

import (
	"github.com/y3933y3933/ghost-card/internal/app"
	"github.com/y3933y3933/ghost-card/internal/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	r := routes.SetRoutes(app)

	r.Run()

}
