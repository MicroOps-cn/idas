package errors

import (
	"strconv"
	"strings"
)

type MultipleError struct {
	errs         []error
	Sep          string
	PluralPrefix string
	Index        bool
	AppendBefore func(err error)
}

func (m *MultipleError) Append(errs ...error) *MultipleError {
	for _, err := range errs {
		if err != nil {
			if m.AppendBefore != nil {
				m.AppendBefore(err)
			}
			m.errs = append(m.errs, errs...)
		}
	}
	return m
}

func (m *MultipleError) HasError() bool {
	return len(m.errs) > 0
}

func AppendError(m *MultipleError, errs ...error) *MultipleError {
	if m != nil {
		_ = m.Append(errs...)
	} else if len(errs) > 0 {
		m = &MultipleError{errs: errs}
	}
	return m
}

func (m MultipleError) Error() string {
	ret := make([]string, len(m.errs))

	for idx, err := range m.errs {
		ret[idx] = ""
		if m.Index {
			ret[idx] = strconv.Itoa(idx) + ". "
		}
		ret[idx] += err.Error()
	}
	if len(ret) > 1 {
		return m.PluralPrefix + strings.Join(ret, m.Sep)
	}
	return strings.Join(ret, m.Sep)
}

func (m *MultipleError) Count() int {
	return len(m.errs)
}

func (m *MultipleError) Clear() {
	m.errs = nil
}

func NewMultipleError(errs ...error) *MultipleError {
	return AppendError(nil, errs...)
}

var _ error = MultipleError{}
