package main

import (
	"fmt"

	"github.com/mijara/idek"
)

type HelloInput struct {
	Name string `json:"name"`
}

type HelloOutput struct {
	Message string `json:"message"`
}

func Hello(ctx *idek.Context, input HelloInput) (*HelloOutput, error) {
	return &HelloOutput{
		Message: fmt.Sprintf("Hello, %s!", input.Name),
	}, nil
}

func main() {
	idek.ViewHandler("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}
