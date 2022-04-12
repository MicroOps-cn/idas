package httputil

import (
	"bytes"
	"strconv"
	"strings"
	"time"
)

type Value struct {
	val       string
	dftVal    string
	splitFunc func(v Value) []string
}

type Option func(value *Value)

func Default(dftVal string) Option {
	return func(value *Value) {
		value.dftVal = dftVal
	}
}

func SplitFunc(f func(v Value) []string) Option {
	return func(value *Value) {
		value.splitFunc = f
	}
}

func NewValue(val string, options ...Option) *Value {
	v := &Value{val: val}
	for _, option := range options {
		option(v)
	}
	return v
}

func (v *Value) Set(val string) *Value {
	v.val = val
	return v
}

func (v *Value) Default(dftVal string) *Value {
	v.dftVal = dftVal
	return v
}

func (v Value) Split(seps ...byte) Values {
	val := v.Bytes()
	if len(val) == 0 {
		return nil
	}
	if bytes.HasPrefix(val, []byte{'['}) && bytes.HasSuffix(val, []byte{']'}) {
		val = bytes.TrimSpace(bytes.TrimSuffix(bytes.TrimPrefix(val, []byte{'['}), []byte{']'}))
	}
	if len(val) == 0 {
		return nil
	}
	if len(seps) == 0 {
		if v.splitFunc != nil {
			var vals Values
			for _, vv := range v.splitFunc(v) {
				vals = append(vals, *NewValue(vv))
			}
			return vals
		}
		seps = []byte{','}
	}
	var vals []Value
loop:
	for pos, i := 0, 0; i < len(val); i++ {
		for _, sep := range seps {
			if sep == val[i] {
				vals = append(vals, Value{val: string(val[pos:i])})
				pos = i + 1
				continue loop
			}
		}
		if i == len(val)-1 {
			vals = append(vals, Value{val: string(val[pos:])})
		}
	}
	return vals
}

func (v Value) String() string {
	if len(v.val) > 0 {
		return v.val
	}
	return v.dftVal
}

func (v Value) Strings(seps ...byte) []string {
	return v.Split(seps...).Strings()
}

func (v Value) Int() (int, error) {
	return strconv.Atoi(strings.TrimSpace(v.String()))
}

func (v Value) Int32s(seps ...byte) (vals []int, err error) {
	return v.Split(seps...).Ints()
}

func (v Value) Int64() (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(v.String()), 10, 64)
}

func (v Value) Int64s(seps ...byte) (vals []int64, err error) {
	return v.Split(seps...).Int64s()
}

func (v Value) Float32() (float32, error) {
	float, err := strconv.ParseFloat(strings.TrimSpace(v.String()), 64)
	return float32(float), err
}

func (v Value) Float32s(seps ...byte) (vals []float32, err error) {
	return v.Split(seps...).Float32s()
}

func (v Value) Float64() (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(v.String()), 64)
}

func (v Value) Float64s(seps ...byte) (vals []float64, err error) {
	return v.Split(seps...).Float64s()
}

func (v Value) Bool() (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(v.String()))
}

func (v Value) Duration() (time.Duration, error) {
	s := strings.TrimSpace(v.String())
	if len(s) > 0 && s[0] == '-' {
		duration, err := time.ParseDuration(s[1:])
		return -duration, err
	}
	return time.ParseDuration(s)
}

func (v Value) Durations(seps ...byte) (vals []time.Duration, err error) {
	return v.Split(seps...).Durations()
}

func (v Value) Time(layout string) (time.Time, error) {
	return time.Parse(layout, v.String())
}

func (v Value) Bytes() []byte {
	return []byte(v.String())
}

type Values []Value

func (v Values) Durations() (vals []time.Duration, err error) {
	for _, s := range v {
		f, err := s.Duration()
		if err != nil {
			return nil, err
		}
		vals = append(vals, f)
	}
	return vals, err
}

func (v Values) Float64s() (vals []float64, err error) {
	for _, s := range v {
		f, err := s.Float64()
		if err != nil {
			return nil, err
		}
		vals = append(vals, f)

	}
	return vals, err
}

func (v Values) Float32s() (vals []float32, err error) {
	for _, s := range v {
		f, err := s.Float32()
		if err != nil {
			return nil, err
		}
		vals = append(vals, f)
	}
	return vals, err
}

func (v Values) Int64s() (vals []int64, err error) {
	for _, s := range v {
		f, err := s.Int64()
		if err != nil {
			return nil, err
		}
		vals = append(vals, f)
	}
	return vals, err
}

func (v Values) Ints() (vals []int, err error) {
	for _, s := range v {
		f, err := s.Int()
		if err != nil {
			return nil, err
		}
		vals = append(vals, f)
	}
	return vals, err
}

func (v Values) Strings() (vals []string) {
	for _, s := range v {
		vals = append(vals, s.String())
	}
	return vals
}
