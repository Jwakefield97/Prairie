package http

/*
	This file will contain structs and methods to model/modify an incoming request. This request struct
	will be passed to the corresponding callback function when a route is matched.
*/

// Request - a struct used to model/modify requests
type Request struct {
	Type    string //Request type (GET,POST,PUT,DELETE). Probably need to make this a struct of constants
	Path    string
	Version string            //http version of the request. Most commonly HTTP/1.1
	Headers map[string]string //headers contained in the request
}

// NewRequest - return an initialized Reqest struct
func NewRequest() Request {
	r := Request{}
	r.Type = ""
	r.Path = ""
	r.Version = ""
	r.Headers = map[string]string{}
	return r
}
