package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

/*
	This file will be used to contain a logger struct and methods for it so that messages from the user/server
	can be logged in the proper format and location.
*/

// TIME_LAYOUT - the layout used to format the time of the log entry
const TIME_LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"

// Log - the logger used to log errors, debug, and access information
type Log struct {
	Path string //location of the log files
}

// NewLog - return a new log instance
func NewLog(path string) Log {
	absPath, _ := filepath.Abs(path)
	l := Log{}
	l.Path = absPath
	return l
}

// CreateLogFiles - a function to create the log files if they dont exist
func CreateLogFiles(log *Log) {
	if _, err := os.Stat(log.Path); os.IsNotExist(err) { //create the logs directory if it doesnt exist
		os.Mkdir(log.Path, os.ModePerm)
	}
	os.OpenFile(log.Path+"/access.txt", os.O_RDONLY|os.O_CREATE, 0666) //create access log
	os.OpenFile(log.Path+"/error.txt", os.O_RDONLY|os.O_CREATE, 0666)  //create error log
	os.OpenFile(log.Path+"/debug.txt", os.O_RDONLY|os.O_CREATE, 0666)  //create debug log
}

// appendToFile - a function to append a given message to a file
func appendToFile(message string, filePath string) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	timeFormatted := time.Now().Format(TIME_LAYOUT)
	log.Printf("%s: %s\n", timeFormatted, message)
}

// Error - log an error
func (l Log) Error(message string) {
	logPath := l.Path + "/error.txt"
	appendToFile(message, logPath)
}

// Debug - log a debug message
func (l Log) Debug(message string) {
	logPath := l.Path + "/debug.txt"
	appendToFile(message, logPath)
}

// Access - log an access message
func (l Log) Access(message string) {
	logPath := l.Path + "/access.txt"
	appendToFile(message, logPath)
}
