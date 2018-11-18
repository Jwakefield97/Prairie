package utils

import (
	"prairie/lib/http"
	"strings"
)

/*
	This file will contain structs/functions that deal with cookies. Everything from parsing to adding and
	verifying them.
*/

// ParseCookies - a function to parse the cookies out of a given request header.
func ParseCookies(request *http.Request) {
	cookieStr := request.Headers["Cookie"]

	if strings.TrimSpace(cookieStr) != "" {
		cookies := strings.Split(cookieStr, ";")
		for _, cookie := range cookies { //loop through all of the cookies
			cookieKeyVal := strings.Split(cookie, "=")
			if len(cookieKeyVal) == 2 {
				request.Cookies[strings.TrimSpace(cookieKeyVal[0])] = cookieKeyVal[1]
			}
		}
	}
}

// AddCookie - add a cookie to a response struct
func AddCookie(response *http.Response) {
	//TODO: implment me
}
