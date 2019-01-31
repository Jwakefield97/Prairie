package main

import (
	"fmt"

	"github.com/Jwakefield97/prairie"
)

// Todo - a struct to test nesting structs in a template
type Todo struct {
	Title string
	Done  bool
}

// TodoPageData - a struct to test template rendering and params
type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

type Name struct {
	Name string
}

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
	app.Get("/builtin/", func(routeObj *prairie.RouteObject) {
		routeObj.Response.File = "templates/builtin.html"
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

	app.Get("/examples", func(routeObj *prairie.RouteObject) {
		routeObj.Response.Template = "examples"
		routeObj.Response.TemplateParams = Name{
			Name: routeObj.Request.Cookies["name"],
		}
	})

	app.Get("/template", func(routeObj *prairie.RouteObject) {
		routeObj.Response.Template = "temp"
		routeObj.Response.TemplateParams = TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
	})

	app.Post("/upload", func(routeObj *prairie.RouteObject) {
		routeObj.Response.SetCookie("name", routeObj.Request.Body["name"], 10000)
		routeObj.Response.Html = "Your name is: " + routeObj.Request.Body["name"] + "<br><a href='/examples'>examples page</a>"
		app.Log.Debug("Uploaded name field: " + routeObj.Request.Body["name"])
		app.Log.Access("Post was hit: " + routeObj.Request.Path)
	})

	app.Start()
}
