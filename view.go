package idek

import (
	"encoding/json"
	"errors"
	"github.com/mozillazg/go-httpheader"
	"io"
	"net/http"
)

type ViewHandler[H, I, O any] func(ctx *Context[H], input I) (O, error)

func View[H, I, O any](path string, handler ViewHandler[H, I, O]) {
	http.Handle(path, http.HandlerFunc(handlerWrapper(handler)))
}

func handlerWrapper[H, I, O any](handler ViewHandler[H, I, O]) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		for _, middleware := range middlewaresHandlers {
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

		input := new(I)

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
			encodeOutput(writer, status, errorResponse)
			return
		}

		encodeOutput(writer, http.StatusOK, output)

	}
}

func encodeOutput(writer http.ResponseWriter, statusCode int, output any) {
	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(output); err != nil {
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
