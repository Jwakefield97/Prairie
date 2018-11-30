package utils

import (
	"github.com/Jwakefield97/prairie/lib/http"
	"strconv"
	"strings"
	"fmt"
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

// ParseMultiPartBody - parse the body of the multipart post request
func ParseMultiPartBody(request *http.Request, body []string) map[string]interface{} {
	//fmt.Println("in parse mulit part")
	returnMap := make(map[string]interface{})
	bodyStr := strings.Join(body,"\n")
	fmt.Println(bodyStr[:200])
	bodyRows := strings.Split(bodyStr,request.BoundaryValue)

	for _, field := range bodyRows { //for each field of the post request
		if strings.TrimSpace(field) != "" { //if not last boundary area (last boundary area ends the post request)

			fieldRows := strings.Split(field,"\n")
			isBody := false
			name := ""
			file := http.UploadFile{}

			for index, row := range fieldRows { //for each row of the field
				if isBody {
					//fmt.Println("in parse mulit part body")
					body := field[index:len(field)] //the rest of the body of the field
					if strings.TrimSpace(file.FileName) != "" { //the field is a file
						//fmt.Println("is file being parsed")

						file.Contents = []byte(body)
						returnMap[name] = file
					}else{ //the field is just a normal field
						returnMap[name] = body
					}
				}else{
					if strings.TrimSpace(row) == "" { //space line between headers and body of the field
						isBody = true
					}else{ //parse headers
						header := strings.Split(row, ":") //split field header
						//fmt.Println(header)
						if len(row) == 2 {
							if strings.EqualFold(strings.TrimSpace(header[0]),"content-disposition") {
								contentFields := strings.Split(header[1],";")

								for _, contentField := range contentFields {
									contentFieldKeyVal := strings.Split(contentField,"=")
									if strings.EqualFold(strings.TrimSpace(contentFieldKeyVal[0]),"name") {
										name = contentFieldKeyVal[1]
									} else if strings.EqualFold(strings.TrimSpace(contentFieldKeyVal[0]),"filename"){
										file.FileName = contentFieldKeyVal[1]
									}
								}
							}else if strings.EqualFold(strings.TrimSpace(header[0]),"content-type") {
								file.FileType = strings.TrimSpace(header[1])
							}
						}
					}
				}
			}
		}
	}
	return returnMap
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
			if isBody {
				if(!request.IsMultiPart){
					params := strings.Split(row, "&")
					for _, param := range params {
						paramKeyVal := strings.Split(param, "=")
						if len(paramKeyVal) == 2 {
							request.Body[strings.TrimSpace(paramKeyVal[0])] = paramKeyVal[1]
						}
					}
				}else{
					request.Body = ParseMultiPartBody(&request,rows[index:len(rows)])
				}
			}else{
				if strings.TrimSpace(row) == "" {
					isBody = true //body parsing has begun
				} else {
					//TODO: parse rest of post request here
					header := strings.SplitN(row, ":", 2) //split on only the first occurence of :
					//TODO: make this more robust
					if len(header) == 2 {
						request.Headers[strings.TrimSpace(header[0])] = header[1]
						if strings.EqualFold(strings.TrimSpace(header[0]),"content-type") { //check if multipart form 
							contentType := strings.SplitN(header[1], ";", 2)
							fmt.Println(contentType[0])
							if strings.EqualFold(strings.TrimSpace(contentType[0]),"multipart/form-data") { //get boundary and store it's key/value
								request.IsMultiPart = true
								boundary := strings.SplitN(contentType[1], "=", 2)
								request.BoundaryKey = strings.TrimSpace(boundary[0])
								request.BoundaryValue = strings.TrimSpace(boundary[1])
							}
						}
					}
				}
			}
		}
	}
	ParseCookies(&request)
	return request
}
