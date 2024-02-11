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
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		request = request.WithContext(context.Background())

		ctx := &Context{
			request: request,
		}

		// Apply middlewares before anything else.
		for _, middleware := range middlewareFuncs {
			if err := middleware(ctx); err != nil {
				encodeError(writer, http.StatusBadRequest, err)
				return
			}
		}

		// This input will concentrate everything, URL Params, Query Params and Body.
		input := new(I)

		// Decode URL Path Params (ex. /hello/:name)
		if err := decoder.Decode(input, transformParams(params)); err != nil {
			encodeError(writer, http.StatusBadRequest, err)
			return
		}

		// Decode query params (ex. ...?query=Hello)
		if err := decoder.Decode(input, request.URL.Query()); err != nil {
			encodeError(writer, http.StatusBadRequest, err)
			return
		}

		// Decode body.
		defer request.Body.Close()
		if err := json.NewDecoder(request.Body).Decode(input); err != nil {
			if !errors.Is(err, io.EOF) {
				encodeError(writer, http.StatusBadRequest, err)
				return
			}
		}

		output, err := handler(ctx, *input)
		if err != nil {
			errorResponse := errorHandler(err)
			if errorResponse == nil {
				errorResponse = &ErrorResponse{
					Error: err.Error(),
				}
			}

			handleOutput(ctx, writer, errorResponse.GetStatus(), errorResponse)
			return
		}

		handleOutput(ctx, writer, http.StatusOK, output)
	}
}

func handleOutput(ctx *Context, writer http.ResponseWriter, statusCode int, output any) {
	for _, onFinishFunc := range ctx.config.onFinishFuncs {
		onFinishFunc(ctx, statusCode, output)
	}

	encoder := json.NewEncoder(writer)
	if ctx.config.pretty {
		encoder.SetIndent("", "  ")
	}

	writer.WriteHeader(statusCode)
	if err := encoder.Encode(output); err != nil {
		encodeError(writer, http.StatusInternalServerError, err)
		return
	}
}

func encodeError(writer http.ResponseWriter, statusCode int, err error) {
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(ErrorResponse{
		Error: err.Error(),
	})
}

func transformParams(params httprouter.Params) url.Values {
	paramValues := url.Values{}
	for _, param := range params {
		paramValues.Set(param.Key, param.Value)
	}
	return paramValues
}
