package main

import (
	"fmt"
	"os"
	"time"
)

type ErrorHandler struct {
	actualTime   time.Time
	errorMessage string
}

// Goroutine: Generate then put error in pipe
func generateError(err error) {
	defer timeTrack(time.Now(), "GENERATE_ERROR")
	var myError = ErrorHandler{
		actualTime:   time.Now(),
		errorMessage: err.Error(),
	}
	if ShellConfig.Error {
		fmt.Fprintln(term, myError.errorMessage)
	}
	g_stdErrPipe <- myError
}

// Goroutine: Generate then put error string in pipe
func generateErrorString(err string) {
	defer timeTrack(time.Now(), "GENERATE_ERROR_STRING")
	var myError = ErrorHandler{
		actualTime:   time.Now(),
		errorMessage: err,
	}
	if ShellConfig.Error {
		fmt.Fprint(term, err)
	}
	g_stdErrPipe <- myError
}

// Goroutine: Write to stderr
func g_errorPipe() {
	// Get
	for elem := range g_stdErrPipe {

		// Set message content
		errorMessage := fmt.Sprintf("[%02d/%02d/%04d - %02d:%02d:%02d]: %s\n",
			elem.actualTime.Day(), elem.actualTime.Month(), elem.actualTime.Year(),
			elem.actualTime.Hour(), elem.actualTime.Minute(), elem.actualTime.Second(), elem.errorMessage)

		out, errno := os.OpenFile(ShellConfig.ErrorFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
		defer out.Close()

		// Output
		if errno == nil {
			out.Write([]byte(errorMessage)) // Error file
		} else {
			os.Stderr.Write([]byte(errorMessage))
		}
	}
}
