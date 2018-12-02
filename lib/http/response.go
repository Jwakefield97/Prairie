package http

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
	This file will contain structs and methods to model/modify an outgoing responses. This response struct
	will be passed to the corresponding callback function when a route is matched so that the user can modify
	response headers and payload.
*/

//TODO: add struct functions to deal with templates and various operations by the user

// Response - a struct to model/modify responses.
type Response struct {
	Status         int
	Template       string      //name of the template to return
	TemplateParams interface{} //array of parameters to pass to the template
	JSON           []byte      //json to return. I'm pretty sure that golang is marshalled into byte arrays. this might need to be updated later
	File           string      //location of the file to be returned
	Html           string
	Text           string            //plain text to be sent back
	Headers        map[string]string //headers to include in the request
	Cookies        []string
	Payload        []byte
}

// NewResponse - constructor for Reponse struct
func NewResponse() Response {
	r := Response{}
	r.Status = HTTP_OK
	r.Template = ""
	r.Html = ""
	r.JSON = make([]byte, 0)
	r.File = ""
	r.Text = ""
	r.Headers = map[string]string{}
	r.Cookies = []string{}
	r.Payload = make([]byte, 0)
	return r
}

// SetCookie - add a cookie to set in a response struct
func (r *Response) SetCookie(key string, val string, seconds int) {
	loc := time.FixedZone("GMT", 0)
	expireTime := time.Now().Add(time.Second * time.Duration(seconds)).In(loc).Format(time.RFC1123)
	cookie := key + "=" + val + "; Expires=" + expireTime + "; HttpOnly; Path=/"
	r.Cookies = append(r.Cookies, cookie)
}

// InvalidateCookie - a function used to invalidate a cookie by setting its date to a date in the past
func (r *Response) InvalidateCookie(key string, val string) {
	cookie := key + "=" + val + "; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Path=/"
	r.Cookies = append(r.Cookies, cookie)
}

// GetDefaultResponse - get the default response struct with preset headers
func GetDefaultResponse() Response {
	r := NewResponse()
	//TODO: add default headers
	return r
}

// GetErrorMessage - returns an appropriate http error response with custom message if supplied
func GetErrorMessage(message string, httpStatus int) *Response {
	response := NewResponse()
	response.Headers["Date"] = time.Now().Format(time.RFC1123)
	response.Headers["Connection"] = "close"
	response.Headers["Server"] = "Prairie"
	response.Headers["Accept-Ranges"] = "bytes"
	response.Headers["Content-Type"] = "text/html"
	response.Status = httpStatus

	response.Payload = []byte(message)

	response.Headers["Content-Length"] = strconv.Itoa(len(response.Payload))

	return &response
}

// FormHTTPResponse - a function to form the actual http response
func FormHTTPResponse(response *Response, templatePath string, canGzip bool) []byte {
	message := make([]byte, 0)
	response.Headers["Date"] = time.Now().Format(time.RFC1123)
	response.Headers["Connection"] = "close"
	response.Headers["Server"] = "Prairie"
	response.Headers["Accept-Ranges"] = "bytes"

	if strings.TrimSpace(response.Html) != "" { //if html response
		response.Payload = []byte(response.Html)
		response.Headers["Content-Type"] = "text/html"

	} else if strings.TrimSpace(string(response.JSON)) != "" { //if json response
		response.Payload = response.JSON
		response.Headers["Content-Type"] = "application/json"

	} else if strings.TrimSpace(response.Text) != "" { //if plain text response
		response.Payload = []byte(response.Text)
		response.Headers["Content-Type"] = "text/plain"

	} else if strings.TrimSpace(response.Template) != "" { //if template response
		absPath, _ := filepath.Abs(templatePath) //get absolute path to templates

		response.Headers["Content-Type"] = "text/html"
		tmpl, _ := template.ParseFiles(absPath + "/" + response.Template + ".p") //parse the template

		var tempBuf bytes.Buffer                                                //buffer to temporarily store template
		if err := tmpl.Execute(&tempBuf, response.TemplateParams); err != nil { //give invalid response if error occurs
			fmt.Println(err)
		}

		response.Payload = []byte(tempBuf.String()) //set the payload of the response

	} else if strings.TrimSpace(response.File) != "" { //if file response
		file := getFile(response.File)
		response.Payload = file.Bytes
		if file.Error == nil && file.Info != nil {
			response.Headers["Last-Modified"] = file.Info.ModTime().Format(time.RFC1123)
			contentType, err := getFileContentType(file.Info.Name())
			if err == nil {
				response.Headers["Content-Type"] = contentType
			} else {
				response = GetErrorMessage("404 not found", HTTP_NOT_FOUND) // if something was wrong with the file name
			}
		} else {
			response = GetErrorMessage("404 not found", HTTP_NOT_FOUND) // file was not found
		}
	}
	if canGzip {
		response.Headers["Content-Encoding"] = "gzip"
		GzipResponseBody(response)
	}
	response.Headers["Content-Length"] = strconv.Itoa(len(response.Payload))
	message = ResponseToBytes(response)
	return message
}

// ResponseToBytes - convert a Reponse struct to a byte array
func ResponseToBytes(response *Response) []byte {
	message := make([]byte, 0)

	message = append(message, []byte("HTTP/1.1 "+strconv.Itoa(response.Status)+" \n")...) //start with status line
	for k, v := range response.Headers {                                                  //append headers
		header := k + ": " + v + "\n"
		message = append(message, []byte(header)...)
	}
	for _, cookie := range response.Cookies { //append headers
		header := "Set-Cookie: " + cookie + "\n"
		message = append(message, []byte(header)...)
	}
	message = append(message, []byte("\n")...) //newline between header and body
	message = append(message, response.Payload...)
	return message
}

// GzipResponseBody - gzip the reponse body of the request
func GzipResponseBody(response *Response) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	defer writer.Close()

	writer.Write(response.Payload)
	response.Payload = buffer.Bytes()

}

// FileStruct - a struct to hold the file data and information about the file
type FileStruct struct {
	Info  os.FileInfo
	Bytes []byte
	Error error
}

// getFileContentType - a function to get the content type of a given file
func getFileContentType(fileName string) (string, error) {
	returnString := ""
	fileTypesArr := strings.Split(fileName, ".")
	if len(fileTypesArr) >= 2 {
		fileType := fileTypesArr[len(fileTypesArr)-1]
		switch fileType {
		case "html":
			returnString = "text/html"
		case "css":
			returnString = "text/css"
		case "js":
			returnString = "application/javascript"
		case "png":
			returnString = "image/png"
		case "jpeg":
			returnString = "image/jpeg"
		case "gif":
			returnString = "audio/mpeg"
		case "mpeg":
			returnString = "audio/mpeg"
		case "json":
			returnString = "application/json"
		case "ico":
			returnString = "image/x-icon"
		default:
			returnString = "text/plain"
		}
	} else {
		return returnString, errors.New("incorrect file name format")
	}
	return returnString, nil
}

//TODO: move this to the proper file
func getFile(name string) FileStruct {
	result := FileStruct{}
	absPath, err := filepath.Abs(name)
	if err != nil {
		result.Error = err
	}
	info, err := os.Stat(absPath) //check if file exists
	//check the error to make sure it is os.IsNotExist(err)
	if err == nil {
		result.Info = info
		// Open file for reading
		file, err := os.Open(absPath)
		defer file.Close()
		if err != nil {
			result.Error = err
		} else {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				result.Error = err
			}
			result.Bytes = data
		}
	} else {
		result.Error = err
	}
	//TODO: change this to have a second return var for errors
	return result //return empty byte array if not found
}
