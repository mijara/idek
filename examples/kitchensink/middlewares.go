package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mijara/idek"
)

func CustomRequestDecoder(ctx *idek.Context, req *http.Request, input any) error {
	return idek.DefaultRequestDecode(ctx, req, input)
}

func CustomResponseEncoder(ctx *idek.Context, rw http.ResponseWriter, res *idek.Response) error {
	return idek.DefaultResponseEncoder(ctx, rw, res)
}

// Example of how to create a gzip response encoder instead of the default JSON encoder.
func CustomGzippedResponseEncoder(ctx *idek.Context, rw http.ResponseWriter, res *idek.Response) error {
	rw.WriteHeader(res.StatusCode)
	rw.Header().Set("Content-Encoding", "gzip")
	rw.Header().Set("Content-Type", "application/json")

	writer := gzip.NewWriter(rw)

	encoder := json.NewEncoder(writer)
	if ctx.EncodingOpts().Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(res)
}

// Example of custom handling of errors returned from endpoints
// Here you can return a statusCode and custom error with data.
func CustomErrorHandler(err error) (int, error) {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest, &CustomError{
			err: err.Error(),
			data: map[string]string{
				"name": "really bad name!",
			},
		}
	}

	return http.StatusInternalServerError, err
}

type CustomError struct {
	err  string
	data map[string]string
}

func (e *CustomError) Error() string {
	return e.err
}

func (e *CustomError) Data() any {
	return e.data
}
