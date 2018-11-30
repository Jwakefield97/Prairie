package http

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie   <-- for setting a cookie
	loc := time.FixedZone("GMT", 0)
	expireTime := time.Now().Add(time.Second * time.Duration(seconds)).In(loc).Format(time.RFC1123)
	cookie := key + "=" + val + "; Expires=" + expireTime + "; HttpOnly; Path=/"
	r.Cookies = append(r.Cookies, cookie)
}

// InvalidateCookie - a function used to invalidate a cookie by setting its date to a date in the past
func (r *Response) InvalidateCookie(key string, val string) {
	//https://stackoverflow.com/questions/5285940/correct-way-to-delete-cookies-server-side  <-- for deleting a cookie
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
func GetErrorMessage(message string, httpStatus int) {
	//TODO: implement me
}

// FormHTTPResponse - a function to form the actual http response
func FormHTTPResponse(response *Response, templatePath string) []byte {
	message := make([]byte, 0)
	response.Headers["Date"] = time.Now().Format(time.RFC1123)
	response.Headers["Connection"] = "close"
	response.Headers["Server"] = "Prairie"
	response.Headers["Accept-Ranges"] = "bytes"

	if strings.TrimSpace(response.Html) != "" {
		response.Payload = []byte(response.Html)
		response.Headers["Content-Type"] = "text/html"

	} else if strings.TrimSpace(string(response.JSON)) != "" {
		response.Payload = response.JSON
		response.Headers["Content-Type"] = "application/json"

	} else if strings.TrimSpace(response.Text) != "" {
		response.Payload = []byte(response.Text)
		response.Headers["Content-Type"] = "text/plain"

	} else if strings.TrimSpace(response.Template) != "" {
		absPath, _ := filepath.Abs(templatePath) //get absolute path to templates

		response.Headers["Content-Type"] = "text/html"
		tmpl, _ := template.ParseFiles(absPath + "/" + response.Template + ".p") //parse the template

		var tempBuf bytes.Buffer                                                //buffer to temporarily store template
		if err := tmpl.Execute(&tempBuf, response.TemplateParams); err != nil { //give invalid response if error occurs
			fmt.Println(err)
		}

		response.Payload = []byte(tempBuf.String()) //set the payload of the response

	} else if strings.TrimSpace(response.File) != "" {
		file := getFile(response.File)
		response.Payload = file.Bytes
		if file.Info != nil {
			response.Headers["Last-Modified"] = file.Info.ModTime().Format(time.RFC1123)
			fileTypesArr := strings.Split(file.Info.Name(), ".")
			//TODO: check to make sure files have a "."
			fileType := fileTypesArr[len(fileTypesArr)-1]
			switch fileType {
			case "html":
				response.Headers["Content-Type"] = "text/html"
			case "css":
				response.Headers["Content-Type"] = "text/css"
			case "js":
				response.Headers["Content-Type"] = "application/javascript"
			case "png":
				response.Headers["Content-Type"] = "image/png"
			case "jpeg":
				response.Headers["Content-Type"] = "image/jpeg"
			case "gif":
				response.Headers["Content-Type"] = "image/gif"
			case "mpeg":
				response.Headers["Content-Type"] = "audio/mpeg"
			case "json":
				response.Headers["Content-Type"] = "application/json"
			default:
				response.Headers["Content-Type"] = "text/plain" //default to plain text if no file type matches
			}
		}

	}
	response.Headers["Content-Length"] = strconv.Itoa(len(response.Payload))

	message = append(message, []byte("HTTP/1.1 200 \n")...) //start with status line
	for k, v := range response.Headers {                    //append headers
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

// FileStruct - a struct to hold the file data and information about the file
type FileStruct struct {
	Info  os.FileInfo
	Bytes []byte
}

//TODO: move this to the proper file
func getFile(name string) FileStruct {
	absPath, _ := filepath.Abs(name)
	result := FileStruct{}
	info, err := os.Stat(absPath) //check if file exists
	//check the error to make sure it is os.IsNotExist(err)
	if err == nil {
		result.Info = info
		// Open file for reading
		file, err := os.Open(absPath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		} else {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatal(err)
			}
			result.Bytes = data
		}
	}
	//TODO: change this to have a second return var for errors
	return result //return empty byte array if not found
}
