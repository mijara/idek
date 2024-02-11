package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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
	time.Sleep(time.Second * 10)

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

	idek.Middleware(idek.SlogMiddleware)
	idek.Middleware(idek.PrettyMiddleware)
	idek.ErrorHandler(ErrorHandler)

	idek.ViewHandler("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}

func ErrorHandler(err error) *idek.ErrorResponse {
	if errors.Is(err, ErrBadRequest) {
		return &idek.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
	}

	// InternalServerError by default.
	return nil
}
