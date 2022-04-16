package api

import "fmt"

type ResponseError struct {
	Code    interface{}
	Message string
}

func NewResponseError(code interface{}, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s[%v]", e.Message, e.Code)
}
