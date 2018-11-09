package prairie

/*
	This is the entry point to the framework. All helper functions/libraries are placed in the folder ./lib.
	MAKE SURE THIS PROJECT IS LOCATED IN YOUR SRC FOLDER OF GO PATH UNDER THE FOLDER "prairie"

	TODO: implement the actual server loop. I think this will be a really good resource: https://golang.org/pkg/net/#example_Listener
	TODO: implement function handleRequest (KEEP PRIVATE)
	TODO: add the Request and Response structs as parameters to RequestCallback
	TODO: add a Response struct to NewPrairieInstance so the user can set a default template for reponses. If on is not passed in
	then a default Reponse struct should be created that autofills certain headers.

	TODO: add an in memory session store.
	TODO: pass the session data structure to the RequestCallback function to be used in the request.
	TODO: add authentication filter to routes. Use a function chaining pattern (like https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/).
	With the chaining style it would look like app.Get("/admin",callBack).isAuthenticated(). Authenticate based on session vars that the user sets.

	TODO: add lib for dealing with JSON: https://golang.org/pkg/compress/gzip/
	TODO: add lib to gzip reponse bodies
	TODO: add template rendering: https://gowebexamples.com/templates/
*/

import (
	"fmt"
	"io"
	"log"
	"net"
	"prairie/lib/http"
	//"prairie/lib/utils"
)

// RouteObject - the object passed to the router methods that holds the request and response.
//TODO: add pointer to session object
type RouteObject struct {
	Request  http.Request
	Response http.Response
}

// RequestCallback - a callback function passed to the Get or Post functions to be called when a url mapping is mapped.
type RequestCallback func(routeObj *RouteObject)

// Prairie - the server struct to act as an interface to the framework.
type Prairie struct {
	ip           string
	port         int
	getMappings  map[string]RequestCallback //all get and post request mappings
	postMappings map[string]RequestCallback //all get and post request mappings
}

// Get - a function for adding a get request mapping to the server.
func (p Prairie) Get(url string, callback RequestCallback) {
	p.getMappings[url] = callback
}

// Post - a function for adding a post request mapping to the server.
func (p Prairie) Post(url string, callback RequestCallback) {
	p.postMappings[url] = callback
}

// NewPrairieInstance - a funciton to create a new Prairie server instance.
//TODO: add an optional Response object parameter where the user can set defaults for the Resonse object to be sent to the router callbacks
func NewPrairieInstance(ip string, port int) Prairie {
	p := Prairie{ip: ip, port: port}
	p.getMappings = map[string]RequestCallback{} //instantiate maps
	p.postMappings = map[string]RequestCallback{}
	return p
}

// Start - a function used to start the server.
// https://golang.org/pkg/net/#example_Listener
func (p Prairie) Start() {
	fmt.Println("The server is being started")
	//TODO: add server loop
	//TODO: spawn routine to handle request

	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 2000,
	}
	//listen for incoming connections
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close() //defer close of listener to the end of infinite for loop
	for {
		// Wait for a connection.
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}

		//spawn new routine to handle incoming connections
		go handleRequest(p, conn)
	}

}

// This will be the function that is used to handle incoming requests in a new go routine
//TODO: add the TCPConn socket object as a parameter to this function
func handleRequest(p Prairie, conn *net.TCPConn) {
	defer conn.Close()
	//read all of the request bytes
	buf := make([]byte, 10000) // 10KB buffer. most browsers limit requests to 8KB. this needs to be changed to be more dynamic
	_, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
	}

	fmt.Println(string(buf))

	//TODO: parse request

	//TODO: set stay alive if the keep alive header is set

	//TODO: create route object based on the request sent and the response template provided at config to pass to callback. This is just a place holder for Request obj
	// routeObj := RouteObject{
	// 	Request: http.Request{
	// 		Path: "this is a test path",
	// 	},
	// 	Response: http.Response{},
	// }

	//TODO: map request to proper route
	//getKeys := reflect.ValueOf(p.getMappings).MapKeys()   //programmatically get the get request keys from the map
	//postKeys := reflect.ValueOf(p.postMappings).MapKeys() //programmatically get the post request keys from the map

	//TODO: call callback for route and
	//p.getMappings[getKeys[0].String()](&routeObj)   //call the callback of the first mapping the in keys for get requests
	//p.postMappings[postKeys[0].String()](&routeObj) //call the callback of the first mapping the in keys for post requests
	//fmt.Println(routeObj.Request.Path)              //print out the modified path from the route REMOVE ME

	//TODO: process response and send it to client
}
