package main

import (
	"github.com/gcapodicimeli/api-movies/cmd/server/routes"
	"github.com/gcapodicimeli/api-movies/pkg/db"
)

func main() {
	engine, db := db.ConnectDatabase()
	router := routes.NewRouter(engine, db)
	router.MapRoutes()

	engine.Run(":8080")
}
