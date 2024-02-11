package idek

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type RequestDecoder func(*Context, *http.Request, any) error

func DefaultRequestDecode(ctx *Context, req *http.Request, input any) error {
	// Decode URL Path Params (ex. /hello/:name)
	if err := decoder.Decode(input, transformParams(ctx.Params())); err != nil {
		return err
	}

	// Decode query params (ex. ...?query=Hello)
	if err := decoder.Decode(input, req.URL.Query()); err != nil {
		return err
	}

	// Decode body.
	if err := json.NewDecoder(req.Body).Decode(input); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}

type ResponseEncoder func(*Context, http.ResponseWriter, *Response) error

type EncoderOptions struct {
	Pretty bool
}

func DefaultResponseEncoder(ctx *Context, rw http.ResponseWriter, res *Response) error {
	rw.WriteHeader(res.StatusCode)

	encoder := json.NewEncoder(rw)
	if ctx.EncodingOpts().Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(res)
}
