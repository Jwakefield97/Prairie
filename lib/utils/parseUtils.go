package utils

import (
	"github.com/Jwakefield97/prairie/lib/http"
	"strconv"
	"strings"
)

// parseQueryParamters - used to parse parameters from the path
func parseQueryParamters(request *http.Request) {
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

// GetContentLength - a function to get only the header/content length from the request string
func GetContentLength(requestStr string) (int, int) {
	contentLen := 0
	header := []string{}
	rows := strings.Split(requestStr, "\n")
	for _, row := range rows {
		header := strings.SplitN(row, ":", 2)
		if strings.EqualFold(header[0], "content-length") {
			i, _ := strconv.Atoi(header[1])
			contentLen = i
		} else if row == "" { //once the body is reached break
			break
		}
		header = append(header, row)
	}
	return len(header), contentLen
}

// ParseHTTPRequest - parse an incoming http request and return a Request struct
func ParseHTTPRequest(requestStr string) http.Request {
	request := http.NewRequest()
	rows := strings.Split(requestStr, "\n")
	isBody := false //whether or not the body is being parsed
	for index, row := range rows {
		if index == 0 { //type,path,version line of http request
			statusStr := strings.Fields(row) //split on whitespace
			//TODO: make this more robust
			if len(statusStr) == 3 {
				request.Type = statusStr[0]
				request.FullPath = statusStr[1]
				request.Version = statusStr[2]
				parseQueryParamters(&request) //parse the parameters out of the path
			}
		} else {
			if row == "" {
				isBody = true //body parsing has begun
			} else if isBody { //parse body params
				params := strings.Split(row, "&")
				for _, param := range params {
					paramKeyVal := strings.Split(param, "=")
					if len(paramKeyVal) == 2 {
						request.Body[strings.TrimSpace(paramKeyVal[0])] = paramKeyVal[1]
					}
				}
			} else {
				//TODO: parse rest of post request here
				header := strings.SplitN(row, ":", 2) //split on only the first occurence of :
				//TODO: make this more robust
				if len(header) == 2 {
					request.Headers[strings.TrimSpace(header[0])] = header[1]
				}
			}
		}
	}
	ParseCookies(&request)
	return request
}
