package utils

import (
	"prairie/lib/http"
	"strings"
)

// parsePathParamters - used to parse parameters from the path
func parsePathParamters(request *http.Request) {
	pathArr := strings.Split(request.FullPath, "?")
	request.Path = pathArr[0] //set Path this should always exist
	if len(pathArr) > 1 {     //parameters are present
		parameters := strings.Split(pathArr[1], "&") //put all parameters in an array
		for _, param := range parameters {
			paramArr := strings.Split(param, "=") //split param into key and value
			if len(paramArr) == 2 {               //if exactly one key val pair
				request.Parameters[paramArr[0]] = paramArr[1]
			}
		}

	}
}

// ParseHTTPRequest - parse an incoming http request and return a Request struct
func ParseHTTPRequest(requestStr string) http.Request {
	request := http.NewRequest()
	rows := strings.Split(requestStr, "\n")
	for index, row := range rows {
		if index == 0 { //type,path,version line of http request
			statusStr := strings.Fields(row) //split on whitespace
			//TODO: make this more robust
			if len(statusStr) == 3 {
				request.Type = statusStr[0]
				request.FullPath = statusStr[1]
				request.Version = statusStr[2]
				parsePathParamters(&request) //parse the parameters out of the path
			}
		} else {
			header := strings.Split(row, ":")
			//TODO: make this more robust
			if len(header) == 2 {
				request.Headers[header[0]] = header[1]
			}
		}
	}
	return request
}
