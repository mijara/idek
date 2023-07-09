package main

import (
	"errors"
	"fmt"
	"idek"
	"net/http"
)

var ErrBadRequest = errors.New("bad request")

type Headers struct {
	UserAgent string `header:"user-agent"`
}

type HelloInput struct {
	Pretty bool   `json:"pretty"`
	Name   string `json:"name"`
}

type HelloOutput struct {
	Message   string `json:"message"`
	UserAgent string `json:"user_agent"`
}

func Hello(ctx *idek.Context[Headers], input HelloInput) (HelloOutput, error) {
	ctx.SetPretty(input.Pretty)

	return HelloOutput{
		Message:   fmt.Sprintf("Hello, %s!", input.Name),
		UserAgent: ctx.Headers().UserAgent,
	}, nil
}

func main() {
	idek.RequestHandler(ValidateHeadersRequestHandler)
	idek.ErrorHandler(ErrorHandler)

	idek.ViewHandler("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}

func ValidateHeadersRequestHandler(request *http.Request) error {
	if xApiToken := request.Header.Get("user-agent"); xApiToken == "" {
		return fmt.Errorf("where is your user-agent!?: %w", ErrBadRequest)
	}

	return nil
}

func ErrorHandler(err error) (int, idek.ErrorResponse) {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest, idek.ErrorResponse{
			Error: err.Error(),
		}
	}

	return http.StatusInternalServerError, idek.ErrorResponse{
		Error: err.Error(),
	}
}
