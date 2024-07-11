package main

import (
	"LeagueManager/internal"
	"LeagueManager/internal/infrastructure/config"
	"LeagueManager/internal/infrastructure/router"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	init, _ := internal.Init()
	app := router.Init(init)

	app.Run(":" + port)
}
