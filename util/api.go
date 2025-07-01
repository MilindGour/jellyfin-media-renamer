// Package util contains utility functions accessed
// by the rest of the application.
package util

import (
	"fmt"
	"log"
	"net/http"
)

func HandleAPIError(w http.ResponseWriter, statusCode int, message string, optionalErrorObject error) {
	msg := fmt.Sprintf("[Error] %d %s.", statusCode, message)
	if optionalErrorObject != nil {
		msg += fmt.Sprintf(" More info: %v", optionalErrorObject)
	}

	w.WriteHeader(statusCode)
	log.Println(msg)   // Write to log
	fmt.Fprint(w, msg) // Write to response
	log.Printf("Message length: %d\n", len(msg))
}
