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
