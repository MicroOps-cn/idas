package errors

import (
	"net/http"
	"strconv"
)

type ServerError interface {
	Code() string
	StatusCode() int
	error
}

var _ ServerError = &serverError{}

type serverError struct {
	code   string
	status int
	err    string
}

func (s serverError) Code() string {
	return s.code
}

func (s serverError) StatusCode() int {
	return s.status
}

func (s serverError) Error() string {
	return s.err
}

func NewServerError(status int, err string, code ...string) ServerError {
	var c string
	if len(code) <= 0 {
		c = strconv.Itoa(status)
	} else {
		c = code[0]
	}
	return &serverError{
		code:   c,
		status: status,
		err:    err,
	}
}

var (
	InternalServerError = NewServerError(http.StatusInternalServerError, "Internal server error")
	NotLoginError       = NewServerError(http.StatusInternalServerError, "Not logged in")
	BadRequestError     = NewServerError(http.StatusBadRequest, "Invalid Request")
	ParameterError      = func(msg string) error { return NewServerError(http.StatusBadRequest, "Parameter Error: "+msg) }
	UnauthorizedError   = NewServerError(http.StatusUnauthorized, "Invalid identity information")
	StatusNotFound      = func(name string) ServerError {
		return NewServerError(http.StatusNotFound, name+" Not Found")
	}
)
