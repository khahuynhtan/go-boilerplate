package main

import (
	"myapp/server"
)

func main() {
	app := server.NewServer()
	if err := app.Start(); err != nil {
		app.Echo.Logger.Fatal("Failed to start server:", err)
		return
	}
	app.Echo.Logger.Fatal(app.Echo.Start(":8080"))

}
