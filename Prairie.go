package prairie

/*
	This is the entry point to the framework. All helper functions/libraries are placed in the folder ./lib.
	MAKE SURE THIS PROJECT IS LOCATED IN YOUR SRC FOLDER OF GO PATH UNDER THE FOLDER "prairie"

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
	"prairie/lib/utils"
	"strings"
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
	ip              string
	port            int
	templateDir     string
	resourceDir     string
	getMappings     map[string]RequestCallback //all get and post request mappings
	postMappings    map[string]RequestCallback //all get and post request mappings
	DefaultResponse http.Response
}

// Get - a function for adding a get request mapping to the server.
func (p Prairie) Get(url string, callback RequestCallback) {
	p.getMappings[url] = callback
}

// Post - a function for adding a post request mapping to the server.
func (p Prairie) Post(url string, callback RequestCallback) {
	p.postMappings[url] = callback
}

// SetTemplateDir - a function to set the resource directory
func (p Prairie) SetTemplateDir(dir string) {
	p.templateDir = dir
}

// SetResourceDir - a function to set the static resource directory
func (p Prairie) SetResourceDir(dir string) {
	p.resourceDir = dir
}

// NewPrairieInstance - a funciton to create a new Prairie server instance.
//TODO: add an optional Response object parameter where the user can set defaults for the Resonse object to be sent to the router callbacks
func NewPrairieInstance(ip string, port int) Prairie {
	p := Prairie{ip: ip, port: port}
	p.getMappings = map[string]RequestCallback{} //instantiate maps
	p.postMappings = map[string]RequestCallback{}
	p.DefaultResponse = http.GetDefaultResponse()
	return p
}

// Start - a function used to start the server.
// https://golang.org/pkg/net/#example_Listener
func (p Prairie) Start() {
	fmt.Printf("The server is being started on %s:%d", p.ip, p.port)

	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	addr := net.TCPAddr{
		IP:   net.ParseIP(p.ip),
		Port: p.port,
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
			continue //continue on error to keep the server up
		}

		//spawn new routine to handle incoming connections
		go handleRequest(p, conn)
	}

}

// This will be the function that is used to handle incoming requests in a new go routine
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

	request := utils.ParseHTTPRequest(string(buf))
	fmt.Println(request.Parameters)
	//TODO: set stay alive if the keep alive header is set

	routeObj := RouteObject{
		Request:  request,
		Response: p.DefaultResponse,
	}

	//match routes and call callback
	if strings.EqualFold(request.Type, "get") {
		if callback, ok := p.getMappings[request.Path]; ok { //if mapping was found
			callback(&routeObj)
		}
	} else if strings.EqualFold(request.Type, "post") {
		if callback, ok := p.postMappings[request.Path]; ok {
			callback(&routeObj)
		}
	}

	//TODO: process response and send it to client
	//responseStr := http.FormHTTPResponse(&routeObj.Response)
	//send response back to the client

}
