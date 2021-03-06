package main

import (
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

/*
	This is a test file to implement the framework and test it as is developed.
*/
func main() {
	app := prairie.NewPrairieInstance("localhost", 2000)
	app.ResourceDir = "resources"
	app.TemplateDir = "templates"
	app.SetLogPath("logs")

	app.Get("/temp", func(routeObj *prairie.RouteObject) {
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

	app.Get("/plain", func(routeObj *prairie.RouteObject) {
		routeObj.Response.SetCookie("lastName", "wakefield", 10000)
		routeObj.Response.File = "templates/test.g"
	})

	app.Get("/index", func(routeObj *prairie.RouteObject) {
		val, ok := routeObj.Session.Load("firstKey")
		if ok {
			routeObj.Response.Html = val.(string) + "<b>Hello from the index page</b>"
		} else {
			routeObj.Response.Html = "<b>Hello from the index page</b>"

		}
	})

	app.Get("/file", func(routeObj *prairie.RouteObject) {
		routeObj.Session.Store("firstKey", "my stored value")
		routeObj.Response.File = "templates/test.html"
	})

	app.Post("/upload", func(routeObj *prairie.RouteObject) {
		routeObj.Response.Text = "Your name is: " + routeObj.Request.Body["name"]
		app.Log.Debug("Uploaded name field: " + routeObj.Request.Body["name"])
		app.Log.Access("Post was hit: " + routeObj.Request.Path)
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
