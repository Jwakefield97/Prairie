package main

import (
	"fmt"

	"github.com/Jwakefield97/prairie"
)

/*
	This is a test file to implement the framework and test it as is developed.
*/
func main() {
	fmt.Println("This is the a driver to test the framework.")
	app := prairie.NewPrairieInstance("localhost", 80)
	app.ResourceDir = "resources"
	app.TemplateDir = "templates"
	app.SetLogPath("logs")

	app.Get("/", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "templates/prairie.html"
	})
	app.Get("/prairie/", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "templates/prairie.html"
	})
	app.Get("/http/", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "templates/http.html"
	})
	app.Get("/utils/", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "templates/utils.html"
	})
	app.Get("/favicon.ico", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "resources/images/favicon.ico"
	})

	app.Get("/logs/error", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "logs/error.txt"
	})

	app.Get("/logs/debug", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "logs/debug.txt"
	})

	app.Get("/logs/access", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "logs/access.txt"
	})

	app.Start()
}
