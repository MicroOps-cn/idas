package errors

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type ServerError interface {
	Code() string
	StatusCode() int
	error
}

func NewMultipleServerError(status int, prefix string, code ...string) *MultipleServerError {
	var c string
	if len(code) <= 0 {
		c = strconv.Itoa(status)
	} else {
		c = code[0]
	}
	return &MultipleServerError{
		code:   c,
		status: status,
		errs:   []error{},
		prefix: prefix,
	}
}

type MultipleServerError struct {
	errs   []error
	code   string
	status int
	prefix string
}

func (m MultipleServerError) Code() string {
	return m.code
}

func (m MultipleServerError) StatusCode() int {
	return m.status
}

func (m MultipleServerError) Error() string {
	if len(m.errs) > 0 {
		if len(m.errs) == 1 {
			return m.errs[0].Error()
		} else {
			var errs []string
			for _, err := range m.errs {
				errs = append(errs, err.Error())
			}
			return m.prefix + strings.Join(errs, ",")
		}
	}
	return ""
}
func (m MultipleServerError) HasError() bool {
	return len(m.errs) > 0
}

func (m MultipleServerError) Append(err error) {
	m.errs = append(m.errs, err)
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
	NotFoundError = NewServerError(404, "record not found")
)

func IsNotFount(err error) bool {
	if err == NotFoundError || err == gorm.ErrRecordNotFound {
		return true
	} else if e, ok := err.(ServerError); ok && e.StatusCode() == 404 {
		return true
	}
	return false
}
