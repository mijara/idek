package view

import (
	"fmt"
	"idek"
	"idek/cmd/app"
)

type HelloInput struct {
	Name string `json:"name"`
}

type HelloOutput struct {
	Response string `json:"response"`
	Headers  any    `json:"headers"`
}

func Hello(ctx *idek.Context[app.Headers], input HelloInput) (HelloOutput, error) {
	return HelloOutput{
		Response: fmt.Sprintf("Hello, %s", input.Name),
		Headers:  ctx.Headers(),
	}, nil
}
