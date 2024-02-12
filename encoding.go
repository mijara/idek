package idek

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type RequestDecoder func(*Context, *http.Request, any) error

type ResponseEncoder func(*Context, http.ResponseWriter, *Response) error

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

func DefaultResponseEncode(ctx *Context, rw http.ResponseWriter, res *Response) error {
	rw.WriteHeader(res.StatusCode)

	encoder := json.NewEncoder(rw)
	if ctx.EncodingOpts().Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(res)
}

func transformParams(params httprouter.Params) url.Values {
	paramValues := url.Values{}
	for _, param := range params {
		paramValues.Set(param.Key, param.Value)
	}
	return paramValues
}
