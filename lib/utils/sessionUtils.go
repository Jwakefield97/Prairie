package utils

/*
	This file will hold all structs/functions to deal with the server's session. It should use mutex's when 
	dealing with the session data structure because it will accessed/modified accross multiple go routines. 
	Potentially look at using the Map from the sync package (https://golang.org/pkg/sync/).
*/