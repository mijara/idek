package idek

import (
	"encoding/json"
)

type ErrorData interface {
	Data() any
}

type Response struct {
	StatusCode int
	Message    any   `json:",omitempty"`
	Error      error `json:",omitempty"`
}

type writableResponse struct {
	StatusCode int
	Message    any            `json:",omitempty"`
	Error      *responseError `json:",omitempty"`
}

type responseError struct {
	Detail string
	Data   any `json:",omitempty"`
}

func (r *Response) IsError() bool {
	return r.Error != nil
}

func (r *Response) MarshalJSON() ([]byte, error) {
	if r.Error != nil {
		var data any
		if errorData, ok := r.Error.(ErrorData); ok {
			data = errorData.Data()
		}

		return json.Marshal(&writableResponse{
			StatusCode: r.StatusCode,
			Error: &responseError{
				Detail: r.Error.Error(),
				Data:   data,
			},
		})
	}

	return json.Marshal(&writableResponse{
		StatusCode: r.StatusCode,
		Message:    r.Message,
	})
}
