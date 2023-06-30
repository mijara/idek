package idek

type ErrorResponse struct {
	Error string
	Data  any
}

type ErrorHandler func(err error) (int, ErrorResponse)

var errorHandler ErrorHandler

func Error(handler ErrorHandler) {
	errorHandler = handler
}
