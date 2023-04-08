/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ServerError interface {
	Code() string
	StatusCode() int
	json.Marshaler
	fmt.Stringer
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

func (m MultipleServerError) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Error())
}

func (m MultipleServerError) String() string {
	return m.Error()
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
		}
		var errs []string
		for _, err := range m.errs {
			errs = append(errs, err.Error())
		}
		return m.prefix + strings.Join(errs, ",")
	}
	return ""
}

func (m MultipleServerError) HasError() bool {
	return len(m.errs) > 0
}

func (m *MultipleServerError) Append(err error) {
	m.errs = append(m.errs, err)
}

var (
	_ ServerError = &MultipleServerError{}
	_ ServerError = &serverError{}
)

type serverError struct {
	code   string
	status int
	err    error
}

func (s *serverError) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Error())
}

func (s serverError) String() string {
	return s.Error()
}

func (s serverError) Code() string {
	return s.code
}

func (s serverError) StatusCode() int {
	return s.status
}

func (s serverError) Error() string {
	return s.err.Error()
}

func (s *serverError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			if f, ok := s.err.(fmt.Formatter); ok {
				f.Format(state, verb)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(state, s.Error())
	case 'q':
		fmt.Fprintf(state, "%q", s.Error())
	}
}

func WithMessage(err error, msg string) error {
	if e, ok := err.(*serverError); ok {
		return &serverError{
			code:   e.code,
			status: e.status,
			err:    errors.WithMessage(e.err, msg),
		}
	} else if e, ok := err.(ServerError); ok {
		return &serverError{
			code:   e.Code(),
			status: e.StatusCode(),
			err:    errors.WithMessage(err, msg),
		}
	}
	if err == gorm.ErrRecordNotFound {
		return &serverError{
			code:   "404",
			status: http.StatusNotFound,
			err:    errors.WithMessage(err, msg),
		}
	}
	return errors.WithMessage(err, msg)
}

func WithServerError(status int, err error, msg string, code ...string) ServerError {
	var c string
	if len(code) == 0 {
		c = strconv.Itoa(status)
	} else {
		c = code[0]
	}
	return &serverError{
		code:   c,
		status: status,
		err:    errors.WithMessage(err, msg),
	}
}

func NewServerError(status int, msg string, code ...string) ServerError {
	var c string
	if len(code) == 0 {
		c = strconv.Itoa(status)
	} else {
		c = code[0]
	}
	return &serverError{
		code:   c,
		status: status,
		err:    errors.New(msg),
	}
}

var (
	InternalServerError = func() error { return NewServerError(http.StatusInternalServerError, "Internal server error") }

	NotLoginError     = func() error { return NewServerError(http.StatusInternalServerError, "Not logged in") }
	BadRequestError   = func() error { return NewServerError(http.StatusBadRequest, "Invalid Request") }
	ParameterError    = func(msg string) error { return NewServerError(http.StatusBadRequest, "Parameter Error: "+msg) }
	UnauthorizedError = func() error { return NewServerError(http.StatusUnauthorized, "Invalid identity information") }
	StatusNotFound    = func(name string) ServerError { return NewServerError(http.StatusNotFound, name+" Not Found") }
	NotFoundError     = func() error { return NewServerError(http.StatusNotFound, "record not found") }
)

func LackParameterError(name string) error {
	return NewServerError(http.StatusBadRequest, "lack parameter: "+name)
}

func IsNotFount(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	} else if e, ok := err.(ServerError); ok && e.StatusCode() == http.StatusNotFound {
		return true
	}
	return false
}

const (
	CodeUserDisable           = "E0001"
	CodeUserNeedResetPassword = "E0002"
)
