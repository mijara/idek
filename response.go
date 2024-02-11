package idek

import (
	"encoding/json"
)

type ErrorData interface {
	Data() any
}

type Response struct {
	// StatusCode is the http status code of this response.
	StatusCode int

	// Message is the OK result of this response.
	Message any `json:",omitempty"`

	// Error is the error of this response.
	// An error can implement the "Data() any" method to extend
	// this into a struct of {Detail, Data} as response.
	Error error `json:",omitempty"`
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
