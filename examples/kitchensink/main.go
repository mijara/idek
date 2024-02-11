package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mijara/idek"
)

var ErrBadRequest = errors.New("bad request")

type HelloInput struct {
	Name string `json:"name"`
}

type HelloOutput struct {
	Message string `json:"message"`
}

func Hello(ctx *idek.Context, input HelloInput) (*HelloOutput, error) {
	time.Sleep(time.Second * 1)

	if input.Name == "bad" {
		return nil, ErrBadRequest
	}

	return &HelloOutput{
		Message: fmt.Sprintf("Hello, %s!", input.Name),
	}, nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	idek.Middleware(idek.PrettyMiddleware)
	idek.Middleware(idek.RequestDecoderMiddleware(CustomRequestDecoder))
	idek.Middleware(idek.ResponseEncoderMiddleware(CustomResponseEncoder))
	idek.Middleware(idek.ErrorMiddleware(CustomErrorHandler)) // important to place this before logging.
	idek.Middleware(idek.SlogMiddleware)

	idek.ViewHandler("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}
