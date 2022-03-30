package httputil

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var ErrStruct = errors.New("Unmarshal() expects struct input. ")

func MapToURLValues(m map[string]string) (vals url.Values) {
	for name, val := range m {
		vals.Set(name, val)
	}
	return
}

// UnmarshalURLValues url.Values to struct
func UnmarshalURLValues(values url.Values, s interface{}) error {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrStruct
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ErrStruct
	}
	return reflectValueFromTag(values, val)
}

func reflectValueFromTag(values url.Values, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		sv := val.Field(i)
		if !(kt.Name[0] >= 'A' && kt.Name[0] <= 'Z') {
			continue
		}
		jsonTag := kt.Tag.Get("json")

		var (
			jsonName = jsonTag
			extAttr  string
		)

		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			jsonName = jsonTag[:idx]
			extAttr = jsonTag[idx+1:]
		}
		if extAttr == "inline" {
			if err := reflectValueFromTag(values, sv); err != nil {
				return err
			}
			continue
		} else if jsonName == "-" {
			continue
		} else if jsonName == "" {
			jsonName = func(old []byte) string {
				if 'A' < old[0] && old[0] < 'Z' {
					old[0] += 'a' - 'A'
				}
				return string(old)
			}([]byte(kt.Name))
		}

		fmt.Println(jsonName, values, sv.Kind())
		switch sv.Kind() {
		case reflect.Struct:
			if err := reflectValueFromTag(values, sv); err != nil {
				return err
			}
			continue
		default:
			if !values.Has(jsonName) {
				continue
			}
		}
		uv := values.Get(jsonName)
		switch sv.Kind() {
		case reflect.Slice:

		case reflect.String:
			sv.SetString(uv)
		case reflect.Bool:
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return fmt.Errorf("cast bool has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, jsonName)
			}
			sv.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || sv.OverflowUint(n) {
				return fmt.Errorf("cast uint has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, jsonName)
			}
			sv.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || sv.OverflowInt(n) {
				return fmt.Errorf("cast int has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, jsonName)
			}
			sv.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(uv, sv.Type().Bits())
			if err != nil || sv.OverflowFloat(n) {
				return fmt.Errorf("cast float has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, jsonName)
			}
			sv.SetFloat(n)
		default:
			return fmt.Errorf("unsupported type: %v ,val: %v ,query key: %v", sv.Type(), uv, jsonName)
		}
	}
	return nil
}
