package idek

import "context"

type contextConfig struct {
	pretty bool
}

type Context[T any] struct {
	ctx     context.Context
	headers T

	// Configs that views can set to change the output format, representation, etc.
	config contextConfig
}

func (c *Context[T]) Ctx() context.Context {
	return c.ctx
}

func (c *Context[T]) Headers() T {
	return c.headers
}

func (c *Context[T]) SetPretty(pretty bool) {
	c.config.pretty = pretty
}
