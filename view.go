package idek

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type ViewHandlerFunc[I, O any] func(ctx *Context, input I) (O, error)

func ViewHandler[I, O any](method, path string, handler ViewHandlerFunc[I, O]) {
	router.Handle(method, path, handlerWrapper(handler))
}

func handlerWrapper[I, O any](handler ViewHandlerFunc[I, O]) httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
		request = request.WithContext(context.Background())

		ctx := &Context{
			request: request,
		}

		// Apply middlewares before anything else.
		for _, middleware := range middlewareFuncs {
			if err := middleware(ctx); err != nil {
				encodeError(w, http.StatusBadRequest, err)
				return
			}
		}

		// This input will concentrate everything, URL Params, Query Params and Body.
		input := new(I)

		// Decode URL Path Params (ex. /hello/:name)
		if err := decoder.Decode(input, transformParams(params)); err != nil {
			encodeError(w, http.StatusBadRequest, err)
			return
		}

		// Decode query params (ex. ...?query=Hello)
		if err := decoder.Decode(input, request.URL.Query()); err != nil {
			encodeError(w, http.StatusBadRequest, err)
			return
		}

		// Decode body.
		defer request.Body.Close()
		if err := json.NewDecoder(request.Body).Decode(input); err != nil {
			if !errors.Is(err, io.EOF) {
				encodeError(w, http.StatusBadRequest, err)
				return
			}
		}

		output, err := handler(ctx, *input)
		if err != nil {
			handleOutput(ctx, w, &Response{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})

			return
		}

		handleOutput(ctx, w, &Response{
			StatusCode: http.StatusOK,
			Message:    output,
		})
	}
}

func handleOutput(ctx *Context, w http.ResponseWriter, res *Response) {
	for _, onFinishFunc := range ctx.config.onFinishFuncs {
		onFinishFunc(ctx, res)
	}

	encoder := json.NewEncoder(w)
	if ctx.config.pretty {
		encoder.SetIndent("", "  ")
	}

	w.WriteHeader(res.StatusCode)
	if err := encoder.Encode(res); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		return
	}
}

func encodeError(writer http.ResponseWriter, statusCode int, err error) {
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(Response{
		StatusCode: statusCode,
		Error:      err,
	})
}

func transformParams(params httprouter.Params) url.Values {
	paramValues := url.Values{}
	for _, param := range params {
		paramValues.Set(param.Key, param.Value)
	}
	return paramValues
}
