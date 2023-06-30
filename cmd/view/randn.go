package view

import (
	"idek"
	"idek/cmd/app"
	"math/rand"
)

type RandInput struct {
	N int `json:"n"`
}

type RandOutput struct {
	Response int `json:"response"`
	Headers  any `json:"headers"`
}

func Rand(ctx *idek.Context[app.Headers], input RandInput) (RandOutput, error) {
	return RandOutput{
		Response: rand.Intn(input.N),
		Headers:  ctx.Headers(),
	}, nil
}
