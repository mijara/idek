package xerrors

import (
	"errors"
	"idek"
	"net/http"
)

func HandleError(err error) (int, idek.ErrorResponse) {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest, idek.ErrorResponse{
			Error: err.Error(),
		}
	}

	return http.StatusInternalServerError, idek.ErrorResponse{
		Error: err.Error(),
	}
}
