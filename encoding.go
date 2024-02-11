package idek

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RequestDecoder func(http.ResponseWriter, *http.Request, httprouter.Params, any) error

func DefaultRequestDecode(w http.ResponseWriter, req *http.Request, params httprouter.Params, input any) error {
	// Decode URL Path Params (ex. /hello/:name)
	if err := decoder.Decode(input, transformParams(params)); err != nil {
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
