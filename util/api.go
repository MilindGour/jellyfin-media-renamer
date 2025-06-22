package util

import (
	"fmt"
	"log"
	"net/http"
)

func HandleAPIError(w http.ResponseWriter, statusCode int, message string, optionalErrorObject error) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("[Error] %s. More info: %v", message, optionalErrorObject)
	fmt.Fprint(w, msg)

	//console log
	consolemsg := fmt.Sprintf("[Error] %d %s. More info: %v", statusCode, message, optionalErrorObject)
	log.Println(consolemsg)
}
