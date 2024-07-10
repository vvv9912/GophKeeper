package customErrors

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Err        error
	StatusCode int
	Message    string
}

func NewCustomError(err error, statusCode int, msg string) *CustomError {
	if err == nil {
		err = errors.New("")
	}
	return &CustomError{err, statusCode, msg}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("[Custom Error Code: %v], [custom msg: %v],[Error: %v]", e.StatusCode, e.Message, e.Err.Error())
}
