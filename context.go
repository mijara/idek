package idek

import (
	"context"
	"net/http"
	"net/url"
)

type Context struct {
	request *http.Request
	config  *contextConfig
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
		o(c.config)
	}
}

type ContextOpt func(o *contextConfig)

type contextConfig struct {
	pretty         bool
	onFinishFuncs  []OnFinishFunc
	requestDecoder RequestDecoder
}

func newDefaultConfig() *contextConfig {
	return &contextConfig{
		requestDecoder: DefaultRequestDecode,
	}
}

func WithPretty(value bool) ContextOpt {
	return func(o *contextConfig) {
		o.pretty = value
	}
}

func WithRequestDecoder(requestDecoder RequestDecoder) ContextOpt {
	return func(o *contextConfig) {
		o.requestDecoder = requestDecoder
	}
}

func WithOnFinish(onFinishFunc OnFinishFunc) ContextOpt {
	return func(o *contextConfig) {
		o.onFinishFuncs = append(o.onFinishFuncs, onFinishFunc)
	}
}
