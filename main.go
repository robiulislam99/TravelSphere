package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "github.com/robiulislam99/TravelSphere/routers"
	"github.com/robiulislam99/TravelSphere/services"
)

func main() {
	services.Init()
	// Set template directory
	web.SetViewsPath("views")

	// Set static file path
	web.SetStaticPath("/static", "static")

	web.Run()
}