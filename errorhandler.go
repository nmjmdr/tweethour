package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tweethour"
)

type ErrorResponse struct {
	ErrorMessage      string `json:"error-message"`
	TwitterStatusCode string `json:"twitter-api-response-status"`
}

func newErrorResponse(message string, code string) *ErrorResponse {
	e := new(ErrorResponse)
	e.ErrorMessage = message
	e.TwitterStatusCode = code
	return e
}

func handleError(w http.ResponseWriter, err tweethour.Error) {

	switch {
	case err.Type() == tweethour.ErrorStatus:
		statusErr, ok := err.(*tweethour.StatusError)
		if !ok {
			fmt.Println("Unable to cast to status error, will just set the error message")
			e := newErrorResponse(err.Error(), "")
			writeError(w, e, NoStatusCode)
		} else {
			e := newErrorResponse(statusErr.Error(), statusErr.StatusText)
			writeError(w, e, statusErr.StatusCode)
		}

	default:
		e := newErrorResponse(err.Error(), "")
		writeError(w, e, NoStatusCode)
	}
}

func writeError(w http.ResponseWriter, e *ErrorResponse, code int) {

	if code != 0 {
		w.WriteHeader(code)
	}

	b, err := json.Marshal((*e))

	if err != nil {
		fmt.Fprintf(w, "Encountered an error, but failed to marshal it")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(b))
}
