package idek

import "net/http"

type ErrorHandlerFunc func(error) *ErrorResponse

var errorHandler ErrorHandlerFunc

func ErrorHandler(handler ErrorHandlerFunc) {
	errorHandler = handler
}

type ErrorResponse struct {
	Status int `json:"-"`
	Error  string
	Data   any
}

func (r *ErrorResponse) GetStatus() int {
	if r.Status == 0 {
		return http.StatusInternalServerError
	}
	return r.Status
}
