package prairie

/*
	This is the entry point to the framework. All helper functions/libraries are placed in the folder ./lib.
	MAKE SURE THIS PROJECT IS LOCATED IN YOUR SRC FOLDER OF GO PATH UNDER THE FOLDER "prairie"

	TODO: add authentication filter to routes. Use a function chaining pattern (like https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/).
	With the chaining style it would look like app.Get("/admin",callBack).isAuthenticated(). Authenticate based on session vars that the user sets.

	TODO: add lib to gzip reponse bodies https://golang.org/pkg/compress/gzip/
*/

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/Jwakefield97/prairie/lib/http"
	"github.com/Jwakefield97/prairie/lib/utils"
)

// BufferSize - the size of the buffer to receive from the socket
const BufferSize = 10000

// Session - the session store to be accessed through routes https://golang.org/pkg/sync/
var Session sync.Map

// RouteObject - the object passed to the router methods that holds the request and response.
type RouteObject struct {
	Request  http.Request
	Response http.Response
	Session  *sync.Map
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
func NewPrairieInstance(ip string, port int) Prairie {
	p := Prairie{ip: ip, port: port}
	p.getMappings = map[string]RequestCallback{} //instantiate maps
	p.postMappings = map[string]RequestCallback{}
	p.DefaultResponse = http.GetDefaultResponse()
	return p
}

// Start - a function used to start the server.
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
	//fmt.Println(requestStr)
	request := utils.ParseHTTPRequest(requestStr)
	//TODO: set stay alive if the keep alive header is set

	//fmt.Println("\n*******************" + request.Cookies["lastName"] + "*********************\n")

	routeObj := RouteObject{
		Request:  request,
		Response: p.DefaultResponse,
		Session:  &Session,
	}
	responseMsg := make([]byte, 0)

	//match routes and call callback
	if strings.EqualFold(request.Type, "get") {
		if callback, ok := p.getMappings[request.Path]; ok { //if mapping was found
			callback(&routeObj)
			responseMsg = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir)
		} else {
			//try to find static resource if not matched by route
			if strings.HasPrefix(request.Path[1:], p.ResourceDir) { //if a public resource was requested
				fmt.Println(request.Path)
				routeObj.Response.File = request.Path[1:]
				responseMsg = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir)
			}
		}
	} else if strings.EqualFold(request.Type, "post") {
		if callback, ok := p.postMappings[request.Path]; ok {
			callback(&routeObj)
			responseMsg = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir)
		}
	}

	//fmt.Println(time.Now().Format(time.RFC1123))
	if len(responseMsg) > 0 { //if less than 0 it is an invalid request
		//fmt.Println(string(responseMsg))
		conn.Write(responseMsg)
	}
	responseMsg = nil

}
