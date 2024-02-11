package idek

import (
	"log/slog"
	"strconv"
	"time"
)

var middlewareFuncs []MiddlewareFunc

type MiddlewareFunc func(ctx *Context) error

func Middleware(handler ...MiddlewareFunc) {
	middlewareFuncs = append(middlewareFuncs, handler...)
}

func PrettyMiddleware(ctx *Context) error {
	prettyStr := ctx.Query().Get("pretty")
	shouldPretty, _ := strconv.ParseBool(prettyStr)
	ctx.Configure(WithPretty(shouldPretty))
	return nil
}

func SlogMiddleware(ctx *Context) error {
	start := time.Now()

	ctx.Configure(WithFinish(func(ctx *Context, status int, output any) {
		slog.Debug("http request",
			slog.String("path", ctx.URL().RawPath),
			slog.Int("status", status),
			slog.Duration("elapsed", time.Since(start)),
		)
	}))

	return nil
}
