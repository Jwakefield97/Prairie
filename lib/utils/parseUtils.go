package utils

import (
	"prairie/lib/http"
	"strings"
)

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
				request.Path = statusStr[1]
				request.Version = statusStr[2]
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
