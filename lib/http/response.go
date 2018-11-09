package http

/*
	This file will contain structs and methods to model/modify an outgoing responses. This response struct
	will be passed to the corresponding callback function when a route is matched so that the user can modify
	response headers and payload.
*/

//TODO: add struct functions to deal with templates and various operations by the user

// Response - a struct to model/modify responses.
type Response struct {
	Status         HttpStatus
	Template       string            //name of the template to return
	TemplateParams []interface{}     //array of parameters to pass to the template
	JSON           []byte            //json to return. I'm pretty sure that golang is marshalled into byte arrays. this might need to be updated later
	File           string            //location of the file to be returned
	Headers        map[string]string //headers to include in the request
	Text           string            //plain text to be sent back
}

// NewResponse - constructor for Reponse struct
func NewResponse() Response {
	r := Response{}
	r.Status = HttpStatus{}
	r.Template = ""
	r.TemplateParams = make([]interface{}, 0)
	r.JSON = make([]byte, 0)
	r.File = ""
	r.Headers = map[string]string{}
	r.Text = ""
	return r
}

// GetDefaultResponse - get the default response struct with preset headers
func GetDefaultResponse() Response {
	r := NewResponse()
	//TODO: add default headers
	return r
}
