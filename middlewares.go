package idek

import (
	"log/slog"
	"strconv"
	"time"
)

var middlewareFuncs []MiddlewareFunc

type MiddlewareFunc func(ctx *Context) error

type OnFinishFunc func(ctx *Context, res *Response)

// Middleware adds all middlewares into the life-cycle of endpoint handling.
// Middlewares are executed before any deserialization and endpoint execution.
func Middleware(funcs ...MiddlewareFunc) {
	middlewareFuncs = append(middlewareFuncs, funcs...)
}

func PrettyMiddleware(ctx *Context) error {
	prettyStr := ctx.Query().Get("pretty")
	shouldPretty, _ := strconv.ParseBool(prettyStr)
	ctx.Configure(WithPretty(shouldPretty))
	return nil
}

func SlogMiddleware(ctx *Context) error {
	start := time.Now()

	ctx.Configure(WithOnFinish(func(ctx *Context, res *Response) {
		attrs := []any{
			slog.String("path", ctx.URL().RawPath),
			slog.Int("status", res.StatusCode),
			slog.Int64("elapsed", time.Since(start).Milliseconds()),
		}

		if res.IsError() {
			attrs = append(attrs, slog.String("error", res.Error.Error()))
		}

		slog.Debug("http request", attrs...)
	}))

	return nil
}

type ErrorHandler func(error) (int, error)

func ErrorMiddleware(handler ErrorHandler) MiddlewareFunc {
	return func(ctx *Context) error {
		ctx.Configure(WithOnFinish(func(ctx *Context, res *Response) {
			if res.Error != nil {
				res.StatusCode, res.Error = handler(res.Error)
			}
		}))

		return nil
	}
}

func RequestDecoderMiddleware(requestDecoder RequestDecoder) MiddlewareFunc {
	return func(ctx *Context) error {
		ctx.Configure(WithRequestDecoder(requestDecoder))
		return nil
	}
}

func ResponseEncoderMiddleware(responseEncoder ResponseEncoder) MiddlewareFunc {
	return func(ctx *Context) error {
		ctx.Configure(WithResponseEncoder(responseEncoder))
		return nil
	}
}
