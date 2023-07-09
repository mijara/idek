package idek

import "net/http"

var requestHandlerFuncs []RequestHandlerFunc

type RequestHandlerFunc func(r *http.Request) error

func RequestHandler(handler RequestHandlerFunc) {
	requestHandlerFuncs = append(requestHandlerFuncs, handler)
}
