package util

import (
	"fmt"
	"net/http"
)

func HandleAPIError(w http.ResponseWriter, statusCode int, message string, optionalErrorObject error) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("[Error] %s. More info: %v", message, optionalErrorObject)
	fmt.Fprint(w, msg)
}
