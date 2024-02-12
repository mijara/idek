package main

import (
	"fmt"

	"github.com/mijara/idek"
)

type HelloInput struct {
	Name string
}

type HelloOutput struct {
	Message string
}

func Hello(ctx *idek.Context, input HelloInput) (*HelloOutput, error) {
	return &HelloOutput{
		Message: fmt.Sprintf("Hello, %s!", input.Name),
	}, nil
}

func main() {
	idek.Handle("GET", "/hello/:name", Hello)

	idek.Start(":8080")
}
