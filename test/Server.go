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
	app := prairie.NewPrairieInstance("127.0.0.1", 3000)

	app.Get("/index", func() {
		fmt.Println("***Inside of the index callback***")
	})

	app.Post("/uploads", func() {
		fmt.Println("***Inside of the uploads callback***")
	})

	app.Start()
}
