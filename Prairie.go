package prairie

/*
	This is the entry point to the framework. All helper functions/libraries are placed in the folder ./lib.

	TODO: add lib to gzip reponse bodies https://golang.org/pkg/compress/gzip/
*/

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Jwakefield97/prairie/lib/http"
	"github.com/Jwakefield97/prairie/lib/utils"
)

// BufferSize - the size of the buffer to receive from the socket
const BufferSize = 10000

// KeepAlivePeriod - how long to keep a connection alive
const KeepAlivePeriod = 30

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
	Log             utils.Log //logger for the prairie instance
}

// NewPrairieInstance - a funciton to create a new Prairie server instance.
func NewPrairieInstance(ip string, port int) Prairie {
	p := Prairie{ip: ip, port: port}
	p.getMappings = map[string]RequestCallback{} //instantiate maps
	p.postMappings = map[string]RequestCallback{}
	p.DefaultResponse = http.GetDefaultResponse()
	p.Log = utils.NewLog("logs")
	return p
}

// SetLogPath - set the location of the log files
func (p *Prairie) SetLogPath(path string) {
	absPath, _ := filepath.Abs(path)
	p.Log.Path = absPath
}

// Get - a function for adding a get request mapping to the server.
func (p Prairie) Get(url string, callback RequestCallback) {
	p.getMappings[url] = callback
}

// Post - a function for adding a post request mapping to the server.
func (p Prairie) Post(url string, callback RequestCallback) {
	p.postMappings[url] = callback
}

// Start - a function used to start the server.
func (p Prairie) Start() {
	fmt.Printf("The server is being started on %s:%d", p.ip, p.port)

	utils.CreateLogFiles(&p.Log)

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
	isKeepAlive := false
	for { // loop until keep alive is no longer valid

		//read all of the request bytes
		buf := make([]byte, BufferSize) // 10KB buffer. most browsers limit requests to 8KB. this needs to be changed to be more dynamic
		bytesRead, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				p.Log.Error("Error occured while reading from connection - " + err.Error())
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
					p.Log.Error("Error occured while reading from connection - " + err.Error())
				}
			}
			requestStr = string(append(buf, newBuf...)) //append the new buffer to the existing buffer
		}

		//fmt.Println(requestStr)
		request := utils.ParseHTTPRequest(requestStr)

		if strings.EqualFold(strings.TrimSpace(request.Headers["Connection"]), "keep-alive") { //if connection is keep alive
			isKeepAlive = true
		}

		routeObj := RouteObject{
			Request:  request,
			Response: p.DefaultResponse,
			Session:  &Session,
		}

		//canGzip := strings.Contains(strings.ToLower(request.Headers["Accept-Encoding"]), "gzip") //check whether the message can be gzipped
		canGzip := false
		responseMsg, responseError := getReponseFromPath(&p, &routeObj, canGzip, isKeepAlive)

		if len(responseMsg) <= 0 { //if less than 0 it is an invalid request
			p.Log.Error("An internal server error was encountered")
			responseMsg = http.ResponseToBytes(http.GetErrorMessage("500 Internal Server Error", http.HTTP_INTERNAL_SERVER_ERROR))
		}

		if isKeepAlive && responseError == nil {
			conn.Write(responseMsg) //send message over connection
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(KeepAlivePeriod * time.Second)
		} else {
			conn.Write(responseMsg) //send message over connection
			break
		}

		responseMsg = nil
	}
}

// getReponseFromPath - get the response from the path whether it is found or not
func getReponseFromPath(p *Prairie, routeObj *RouteObject, canGzip bool, isKeepAlive bool) ([]byte, error) {
	responseMsg := make([]byte, 0)
	var returnError error

	//match routes and call callback
	if strings.EqualFold(routeObj.Request.Type, "get") {
		if callback, ok := p.getMappings[routeObj.Request.Path]; ok { //if mapping was found
			callback(routeObj)
			responseMsg, returnError = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir, canGzip, isKeepAlive, KeepAlivePeriod)
		} else {
			//try to find static resource if not matched by route
			if strings.HasPrefix(routeObj.Request.Path[1:], p.ResourceDir) { //if a public resource was requested
				routeObj.Response.File = routeObj.Request.Path[1:]
				responseMsg, returnError = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir, canGzip, isKeepAlive, KeepAlivePeriod)
			} else {
				p.Log.Access("Not Found (GET) - " + routeObj.Request.Path) // log path not found
				responseMsg = http.ResponseToBytes(http.GetErrorMessage("404 Path Not Found", http.HTTP_NOT_FOUND))
				returnError = errors.New("post mapping not found")
			}
		}
	} else if strings.EqualFold(routeObj.Request.Type, "post") {
		if callback, ok := p.postMappings[routeObj.Request.Path]; ok {
			callback(routeObj)
			responseMsg, returnError = http.FormHTTPResponse(&routeObj.Response, p.TemplateDir, canGzip, isKeepAlive, KeepAlivePeriod)
		} else {
			p.Log.Access("Not Found (POST) - " + routeObj.Request.Path) // log path not found
			responseMsg = http.ResponseToBytes(http.GetErrorMessage("404 Path Not Found", http.HTTP_NOT_FOUND))
			returnError = errors.New("post mapping not found")
		}

	} else {
		p.Log.Access("Unknown Request Type (" + routeObj.Request.Type + ") - " + routeObj.Request.Path) // unknown request type
		responseMsg = http.ResponseToBytes(http.GetErrorMessage("404 Path Not Found", http.HTTP_METHOD_NOT_ALLOWED))
		returnError = errors.New("post mapping not found")
	}
	return responseMsg, returnError
}
