package client

import (
	"errors"
	"fmt"
)

type StatusError struct {
	Code int
	Body string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("unexpected status %d: %s", e.Code, e.Body)
}

func NewStatusError(code int, body string) *StatusError {
	return &StatusError{Code: code, Body: body}
}

func IsStatusCode(err error, code int) bool {
	var se *StatusError
	if errors.As(err, &se) {
		return se.Code == code
	}
	return false
}
