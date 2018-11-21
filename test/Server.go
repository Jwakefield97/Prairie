package main

import (
	"fmt"
	"prairie"
)

/*
	This is a test file to implement the framework and test it as is developed.
*/
func main() {
	fmt.Println("This is the a driver to test the framework.")
	app := prairie.NewPrairieInstance("127.0.0.1", 2000)
	app.ResourceDir = "resources"
	app.TemplateDir = "templates"

	app.Get("/index", func(routeObj *prairie.RouteObject) {
		fmt.Println("***Inside of the index callback***")
		routeObj.Response.Html = "<b>Hello from the index page</b>"
	})

	app.Get("/file", func(routeObj *prairie.RouteObject) {
		fmt.Println("***Inside of the index callback***")
		routeObj.Response.File = "templates/test.html"
	})

	app.Post("/uploads", func(routeObj *prairie.RouteObject) {
		fmt.Println("***Inside of the uploads callback***")
		routeObj.Response.Text = "This is plain text"
	})

	app.Start()
}
