package main

import (
	"log"

	"github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
	_ "github.com/robiulislam99/TravelSphere/routers"
	"github.com/robiulislam99/TravelSphere/services"
)

func main() {
	// Load .env FIRST — before anything reads os.Getenv
	if err := godotenv.Load(); err != nil {
		log.Println("[main] WARNING: no .env file found, relying on system environment")
	}
	services.Init()
	// Set template directory
	web.SetViewsPath("views")

	// Set static file path
	web.SetStaticPath("/static", "static")

	web.Run()
}
