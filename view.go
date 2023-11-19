package idek

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/mozillazg/go-httpheader"
	"io"
	"net/http"
	"net/url"
)

type ViewHandlerFunc[H, I, O any] func(ctx *Context[H], input I) (O, error)

func ViewHandler[H, I, O any](method, path string, handler ViewHandlerFunc[H, I, O]) {
	router.Handle(method, path, handlerWrapper(handler))
}

func handlerWrapper[H, I, O any](handler ViewHandlerFunc[H, I, O]) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// Apply middlewares before anything else.
		for _, middleware := range requestHandlerFuncs {
			if err := middleware(request); err != nil {
				encodeError(writer, http.StatusBadRequest, err)
				return
			}
		}

		// Decode Headers.
		headers := new(H)
		if err := httpheader.Decode(request.Header, headers); err != nil {
			encodeError(writer, http.StatusBadRequest, err)
			return
		}

		// This input will concentrate everything, URL Params, Query Params and Body.
		input := new(I)

		// Decode URL Params.
		if err := decoder.Decode(input, transformParams(params)); err != nil {
			encodeError(writer, http.StatusBadRequest, err)
			return
		}

		// Decode query params.
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

		ctx := &Context[H]{
			ctx:     request.Context(),
			headers: *headers,
		}

		output, err := handler(ctx, *input)
		if err != nil {
			status, errorResponse := errorHandler(err)
			encodeOutput(contextConfig{}, writer, status, errorResponse)
			return
		}

		encodeOutput(ctx.config, writer, http.StatusOK, output)
	}
}

func encodeOutput(config contextConfig, writer http.ResponseWriter, statusCode int, output any) {
	encoder := json.NewEncoder(writer)
	if config.pretty {
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
