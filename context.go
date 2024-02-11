package idek

import (
	"context"
	"net/http"
	"net/url"
)

type contextConfig struct {
	pretty        bool
	onFinishFuncs []OnFinishFunc
}

type Context struct {
	request *http.Request

	// Configs that views can set to change the output format, representation, etc.
	config contextConfig
}

func (c *Context) Ctx() context.Context {
	return c.request.Context()
}

func (c *Context) URL() *url.URL {
	return c.request.URL
}

func (c *Context) Header() http.Header {
	return c.request.Header
}

func (c *Context) Query() url.Values {
	return c.request.URL.Query()
}

func (c *Context) Configure(opts ...ContextOpt) {
	for _, o := range opts {
		o(&c.config)
	}
}

type ContextOpt func(o *contextConfig)

func WithPretty(value bool) ContextOpt {
	return func(o *contextConfig) {
		o.pretty = value
	}
}

func WithFinish(onFinishFunc OnFinishFunc) ContextOpt {
	return func(o *contextConfig) {
		o.onFinishFuncs = append(o.onFinishFuncs, onFinishFunc)
	}
}
