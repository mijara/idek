package idek

type ErrorResponse struct {
	Error string
	Data  any
}

type ErrorHandlerFunc func(error) (int, ErrorResponse)

var errorHandler ErrorHandlerFunc

func ErrorHandler(handler ErrorHandlerFunc) {
	errorHandler = handler
}
