package idek

import "context"

type Context[T any] struct {
	ctx     context.Context
	headers T
}

func (c *Context[T]) Ctx() context.Context {
	return c.ctx
}

func (c *Context[T]) Headers() T {
	return c.headers
}
