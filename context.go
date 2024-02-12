package idek

import (
	"context"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type Context struct {
	request *http.Request
	params  httprouter.Params // To be deprecated in go1.22
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

func (c *Context) Params() httprouter.Params {
	return c.params
}

func (c *Context) Configure(opts ...ContextOpt) {
	for _, o := range opts {
		o(c.config)
	}
}

type EncoderOptions struct {
	Pretty bool
}

func (c *Context) EncodingOpts() EncoderOptions {
	return EncoderOptions{
		Pretty: c.config.pretty,
	}
}

type ContextOpt func(o *contextConfig)

type contextConfig struct {
	pretty          bool
	onFinishFuncs   []OnFinishFunc
	requestDecoder  RequestDecoder
	responseEncoder ResponseEncoder
}

func newDefaultConfig() *contextConfig {
	return &contextConfig{
		requestDecoder:  DefaultRequestDecode,
		responseEncoder: DefaultResponseEncoder,
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

func WithResponseEncoder(responseEncoder ResponseEncoder) ContextOpt {
	return func(o *contextConfig) {
		o.responseEncoder = responseEncoder
	}
}

func WithOnFinish(onFinishFunc OnFinishFunc) ContextOpt {
	return func(o *contextConfig) {
		o.onFinishFuncs = append(o.onFinishFuncs, onFinishFunc)
	}
}
