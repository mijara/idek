package middleware

import (
	"fmt"
	"idek/cmd/xerrors"
	"net/http"
)

func ValidateHeaders(request *http.Request) error {
	if xApiToken := request.Header.Get("x-api-token"); xApiToken == "" {
		return fmt.Errorf("x-api-token missing from request: %w", xerrors.ErrBadRequest)
	}

	return nil
}
