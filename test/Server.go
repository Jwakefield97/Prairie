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

	app.Get("/index", func(routeObj *prairie.RouteObject) {
		fmt.Println("***Inside of the index callback***")
		fmt.Println(routeObj.Request.Parameters)
	})

	app.Post("/uploads", func(routeObj *prairie.RouteObject) {
		fmt.Println("***Inside of the uploads callback***")
		fmt.Println(routeObj.Request.Path)
		routeObj.Request.Path = "this is the new path from the callback"
	})

	app.Start()
}
