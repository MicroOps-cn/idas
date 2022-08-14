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

package httputil

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ErrStruct = errors.New("Unmarshal() expects struct input. ")

func MapToURLValues(m map[string]string) (vals url.Values) {
	vals = make(url.Values)
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
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return reflectValueFromTag(values, val.Elem())
	}
loop:
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
				if enum, ok := sv.Interface().(protoreflect.Enum); ok {
					if enumVal := enum.Descriptor().Values().ByName(protoreflect.Name(uv)); enumVal != nil {
						sv.SetInt(int64(enumVal.Number()))
						continue
					} else if enumVal = enum.Descriptor().Values().ByName(protoreflect.Name(strings.ToLower(uv))); enumVal != nil {
						sv.SetInt(int64(enumVal.Number()))
						continue
					} else if enumVal = enum.Descriptor().Values().ByName(protoreflect.Name(strings.ToUpper(uv))); enumVal != nil {
						sv.SetInt(int64(enumVal.Number()))
						continue
					} else if enumVal = enum.Descriptor().Values().ByName(protoreflect.Name(strings.ToTitle(uv))); enumVal != nil {
						sv.SetInt(int64(enumVal.Number()))
						continue
					}
				}

				if protoTag := kt.Tag.Get("protobuf"); protoTag != "" {
					var typeName string
					for _, s := range strings.Split(protoTag, ",") {
						if strings.HasPrefix(s, "enum=") {
							typeName = s[5:]
							break
						}
					}
					if len(typeName) != 0 {
						enumMap := proto.EnumValueMap(typeName)
						for v, x := range enumMap {
							if v == uv {
								sv.SetInt(int64(x))
								continue loop
							}
						}
					}
				}
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
