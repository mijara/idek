package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
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
	idek.Middleware(idek.ErrorMiddleware(ErrorHandler)) // Important to place this before logging
	idek.Middleware(idek.SlogMiddleware)

	idek.ViewHandler("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}

func CustomRequestDecoder(w http.ResponseWriter, req *http.Request, params httprouter.Params, input any) error {
	return idek.DefaultRequestDecode(w, req, params, input)
}

func ErrorHandler(err error) (int, error) {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest, &MyCustomError{
			err: err.Error(),
			data: map[string]string{
				"name": "really bad name!",
			},
		}
	}

	return http.StatusInternalServerError, err
}

type MyCustomError struct {
	err  string
	data map[string]string
}

func (e *MyCustomError) Error() string {
	return e.err
}

func (e *MyCustomError) Data() any {
	return e.data
}
