package http

import (
	"io/ioutil"
	"path/filepath"
)

/*
	This file will contain structs and methods to model/modify an incoming request. This request struct
	will be passed to the corresponding callback function when a route is matched.
*/

// Request - a struct used to model/modify requests
type Request struct {
	Type       string            //Request type (GET,POST,PUT,DELETE). Probably need to make this a struct of constants
	FullPath   string            //full path including paramters
	Path       string            //just the path
	Version    string            //http version of the request. Most commonly HTTP/1.1
	Headers    map[string]string //headers contained in the request
	Parameters map[string]string //parameters from the path
	Cookies    map[string]string //a map of cookies from the http request
	Body       map[string]string //body of a post request
}

// NewRequest - return an initialized Reqest struct
func NewRequest() Request {
	r := Request{}
	r.Type = ""
	r.Path = ""
	r.FullPath = ""
	r.Version = ""
	r.Headers = map[string]string{}
	r.Parameters = map[string]string{}
	r.Cookies = map[string]string{}
	r.Body = map[string]string{}
	return r
}

// UploadFile - a struct to represent the uploaded file from a post request
type UploadFile struct {
	Contents []byte
	FileType string
	FileName string
}

// Save - save a given file from a file upload
func (f UploadFile) Save(location string) {
	absPath, _ := filepath.Abs(location + f.FileName)
	err := ioutil.WriteFile(absPath, f.Contents, 0644)
	if err != nil {
		panic(err) //TODO: change this so the server doesnt crash
	}
}
