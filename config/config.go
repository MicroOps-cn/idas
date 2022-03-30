package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
)

func ref(path string, val reflect.Value) interface{} {
	if len(path) == 0 {
		return val.Interface()
	}
	switch val.Kind() {
	case reflect.Ptr:
		return ref(path, val.Elem())
	}
	if len(path) >= 1 && path[0] == '.' {
		return ref(path[1:], val.Elem())
	}
	typ := val.Type()

	fmt.Println(path, val.Type().Name())
	idx := strings.IndexAny(path, ".[")
	if idx >= 0 {
		var downPath = path[idx:]
		if path[idx] == '.' {
			downPath = path[idx+1:]
		} else if path[idx] == '[' && idx == 0 {
			idx2 := strings.IndexRune(path[idx+1:], ']')
			fmt.Println(typ.Kind())
			if idx2 >= 0 && typ.Kind() == reflect.Slice {
				index, err := strconv.Atoi(path[idx+1 : idx+1+idx2])
				if err != nil {
					return nil
				}
				return ref(path[idx+1+idx2+1:], val.Index(index))
			}
			return nil
		}
		for i := 0; i < val.NumField(); i++ {
			kt := typ.Field(i)
			sv := val.Field(i)
			if !(kt.Name[0] >= 'A' && kt.Name[0] <= 'Z') {
				continue
			}
			var jsonName string
			jsonTag := kt.Tag.Get("json")
			if i1 := strings.Index(jsonTag, ","); i1 >= 0 {
				jsonName = jsonTag[:i1]
			} else {
				for i2, c := range kt.Name {
					if !(c >= 'A' && c <= 'Z') {
						if i2 != 0 {
							jsonName += "_"
						}
						jsonName += string([]int32{c + ('a' - 'A')})
						continue
					}
					jsonName += string([]int32{c})
				}
			}
			fmt.Println(jsonName, path[:idx])
			if jsonName == path[:idx] {
				return ref(downPath, sv)
			}
		}
	}
	return nil
}
func (x *Config) Init(logger log.Logger) error {
	for _, userStorage := range x.Storage.User {
		switch s := userStorage.Source.(type) {
		case *Storage_Ref:
			s.Ref.Storage = ref(s.Ref.Path, reflect.ValueOf(x)).(*Storage)
		}
	}
	return nil
}

func (x *StorageRef) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, b []byte) error {
	return json.Unmarshal(b, &x.Path)
}
func (m *Storage) GetStorageSource() isStorage_Source {
	if m != nil {
		switch s := m.Source.(type) {
		case *Storage_Ref:
			return s.Ref.GetStorage().GetStorageSource()
		default:
			return m.Source
		}
	}
	return nil
}
func (x *MySQLOptions) GetStdMaxConnectionLifeTime() time.Duration {
	if x != nil {
		if duration, err := types.DurationFromProto(x.MaxConnectionLifeTime); err == nil {
			return duration
		}
	}
	return time.Second * 30
}

func (x *MySQLOptions) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, b []byte) error {
	options := NewMySQLOptions()
	x.Charset = options.Charset
	x.Collation = options.Collation
	x.MaxIdleConnections = options.MaxIdleConnections
	x.MaxOpenConnections = options.MaxOpenConnections
	x.MaxConnectionLifeTime = options.MaxConnectionLifeTime
	x.TablePrefix = options.TablePrefix
	return json.Unmarshal(b, x)
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Charset:               "utf8",
		Collation:             "utf8_general_ci",
		MaxIdleConnections:    2,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: types.DurationProto(30 * time.Second),
		TablePrefix:           "t_",
	}
}
