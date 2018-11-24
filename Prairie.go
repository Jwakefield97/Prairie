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

// BufferSize - the size of the buffer to receive from the socket
const BufferSize = 10000

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
	TemplateDir     string
	ResourceDir     string
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
	buf := make([]byte, BufferSize) // 10KB buffer. most browsers limit requests to 8KB. this needs to be changed to be more dynamic
	bytesRead, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
	}
	requestStr := string(buf)
	headerLen, contentLen := utils.GetContentLength(requestStr)

	remainingBytes := contentLen - (bytesRead - headerLen) //calculate remaining bytes to be read

	if remainingBytes > 0 { //if there is more content to pulled from the socket
		newBuf := make([]byte, remainingBytes) //make a new buffer at the size of the existing bytes to pull from socket
		_, err := conn.Read(newBuf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
		}
		requestStr = string(append(buf, newBuf...)) //append the new buffer to the existing buffer
	}

	request := utils.ParseHTTPRequest(requestStr)
	//TODO: set stay alive if the keep alive header is set

	routeObj := RouteObject{
		Request:  request,
		Response: p.DefaultResponse,
	}
	responseMsg := make([]byte, 0)

	//match routes and call callback
	if strings.EqualFold(request.Type, "get") {
		if callback, ok := p.getMappings[request.Path]; ok { //if mapping was found
			callback(&routeObj)
			responseMsg = http.FormHTTPResponse(&routeObj.Response)
		} else {
			//try to find static resource if not matched by route
			if strings.HasPrefix(request.Path[1:], p.ResourceDir) { //if a public resource was requested
				routeObj.Response.File = request.Path[1:]
				responseMsg = http.FormHTTPResponse(&routeObj.Response)
			}
		}
	} else if strings.EqualFold(request.Type, "post") {
		if callback, ok := p.postMappings[request.Path]; ok {
			callback(&routeObj)
			responseMsg = http.FormHTTPResponse(&routeObj.Response)
		}
	}

	//fmt.Println(time.Now().Format(time.RFC1123))
	if len(responseMsg) > 0 { //if less than 0 it is an invalid request
		fmt.Println(string(responseMsg))
		conn.Write(responseMsg)
	}

}
