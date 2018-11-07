package prairie

/*
	This is the entry point to the framework. All helper functions/libraries are placed in the folder ./lib.
	MAKE SURE THIS PROJECT IS LOCATED IN YOUR SRC FOLDER OF GO PATH UNDER THE FOLDER "prairie"

	TODO: implement the actual server loop. I think this will be a really good resource: https://golang.org/pkg/net/#example_Listener
	TODO: implement function handleRequest (KEEP PRIVATE)
	TODO: add the Request and Response structs as parameters to RequestCallback
	TODO: add a Response struct to NewPrairieInstance so the user can set a default template for reponses. If on is not passed in 
	then a default Reponse struct should be created that autofills certain headers.
	TODO: add lib for dealing with JSON: https://golang.org/pkg/compress/gzip/
	TODO: add lib to gzip reponse bodies 
	TODO: add an in memory session store. 
	TODO: add template rendering: https://gowebexamples.com/templates/
*/

import (
	"fmt"
	"reflect"
	//"prairie/lib/http"
	//"prairie/lib/utils"
)

// RequestCallback - a callback function passed to the Get or Post functions to be called when a url mapping is mapped.
type RequestCallback func()

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
	getKeys := reflect.ValueOf(p.getMappings).MapKeys()   //programmatically get the get request keys from the map
	postKeys := reflect.ValueOf(p.postMappings).MapKeys() //programmatically get the post request keys from the map

	p.getMappings[getKeys[0].String()]()   //call the callback of the first mapping the in keys for get requests
	p.postMappings[postKeys[0].String()]() //call the callback of the first mapping the in keys for post requests
}

// This will be the function that is used to handle incoming requests in a new go routine
func handleRequest() {

}
