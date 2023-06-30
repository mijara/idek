package idek

import "net/http"

var middlewaresHandlers []MiddlewareHandler

type MiddlewareHandler func(r *http.Request) error

func Middleware(handler MiddlewareHandler) {
	middlewaresHandlers = append(middlewaresHandlers, handler)
}
