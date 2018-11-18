package main

import (
	"fmt"
	"prairie/lib/utils"
)

/*
	This is a test file to implement the framework and test it as is developed.
*/
func main() {
	// fmt.Println("This is the a driver to test the framework.")
	// app := prairie.NewPrairieInstance("127.0.0.1", 2000)
	// app.SetResourceDir("resources")
	// app.SetTemplateDir("templates")

	// app.Get("/index", func(routeObj *prairie.RouteObject) {
	// 	fmt.Println("***Inside of the index callback***")
	// 	fmt.Println(routeObj.Request.Parameters)
	// })

	// app.Post("/uploads", func(routeObj *prairie.RouteObject) {
	// 	fmt.Println("***Inside of the uploads callback***")
	// 	fmt.Println(routeObj.Request.Path)
	// 	routeObj.Request.Path = "this is the new path from the callback"
	// })

	// app.Start()
	str := `
	Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
	Accept-Encoding: gzip, deflate, br
	Accept-Language: en-US,en;q=0.9
	Cache-Control: max-age=0
	Connection: keep-alive
	Content-Length: 13
	Content-Type: application/x-www-form-urlencoded
	Cookie: name=jake; cart=empty; 
	Origin: http://localhost:2000
	Referer: http://localhost:2000/
	Upgrade-Insecure-Requests: 1
	User-Agent: Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36

	say=Hi&to=Mom
	1=Hi&2=Mom
	3=Hi&4=Mom
	`
	req := utils.ParseHTTPRequest(str)
	fmt.Println(req.Body)
}
