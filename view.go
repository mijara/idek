package idek

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc[I, O any] func(ctx *Context, input I) (O, error)

func Handle[I, O any](method, path string, handler HandlerFunc[I, O]) {
	router.Handle(method, path, handlerWrapper(handler))
}

func handlerWrapper[I, O any](handler HandlerFunc[I, O]) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := &Context{
			request: req.WithContext(context.Background()),
			config:  newDefaultConfig(),
			params:  params,
		}

		// Apply middlewares before anything else.
		for _, middleware := range middlewareFuncs {
			if err := middleware(ctx); err != nil {
				encodeError(w, http.StatusBadRequest, err)
				return
			}
		}

		// Decode endpoint input according to generic type.
		input := new(I)
		if err := ctx.config.requestDecoder(ctx, ctx.request, input); err != nil {
			encodeError(w, http.StatusBadRequest, err)
			return
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

	if err := ctx.config.responseEncoder(ctx, w, res); err != nil {
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
