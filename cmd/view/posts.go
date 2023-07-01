package view

import (
	"idek"
	"idek/cmd/app"
)

type Posts struct {
}

func NewPosts() *Posts {
	return &Posts{}
}

type AllPostsInput struct {
	Limit int
}

type AllPostsOutput struct {
	Count int
}

func (p *Posts) All(ctx *idek.Context[app.Headers], input AllPostsInput) (AllPostsOutput, error) {
	return AllPostsOutput{
		Count: input.Limit,
	}, nil
}
