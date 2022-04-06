package config

import (
	"bytes"
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
	if val.Kind() == reflect.Ptr {
		return ref(path, val.Elem())
	}
	if len(path) == 0 {
		return val.Interface()
	}

	if len(path) >= 1 && path[0] == '.' {
		return ref(path[1:], val.Elem())
	}
	typ := val.Type()

	idx := strings.IndexAny(path, ".[")
	if idx >= 0 {
		downPath := path[idx:]
		if path[idx] == '.' {
			downPath = path[idx+1:]
		} else if path[idx] == '[' && idx == 0 {
			idx2 := strings.IndexRune(path[idx+1:], ']')
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
			if jsonName == path[:idx] {
				return ref(downPath, sv)
			}
		}
	}
	return nil
}

func (x *Storage) findRef(path string, root interface{}) error {
	target := ref(path, reflect.ValueOf(root))
	buf := bytes.Buffer{}
	unmarshaller := jsonpb.Marshaler{}
	switch s := target.(type) {
	case Storage:
		if s.GetRef() != nil {
			return x.findRef(s.GetRef().Path, root)
		} else if err := unmarshaller.Marshal(&buf, &s); err != nil {
			return err
		}
	case *Storage:
		if s.GetRef() != nil {
			return x.findRef(s.GetRef().Path, root)
		} else if err := unmarshaller.Marshal(&buf, s); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown ref: %s(%T)", path, target)
	}
	tmpStorage := new(Storage)
	if err := jsonpb.Unmarshal(&buf, tmpStorage); err != nil {
		return err
	}
	x.Source = tmpStorage.Source
	return nil
}

func (x *Config) Init(logger log.Logger) error {
	for _, userStorage := range append(append(x.Storage.User, x.Storage.App...), x.Storage.Session) {
		switch s := userStorage.Source.(type) {
		case *Storage_Ref:
			if s.Ref.Storage == nil {
				s.Ref.Storage = new(Storage)
			}
			err := s.Ref.Storage.findRef(s.Ref.Path, x)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *StorageRef) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, b []byte) error {
	return json.Unmarshal(b, &x.Path)
}

func (x *Storage) GetStorageSource() isStorage_Source {
	if x != nil {
		switch s := x.Source.(type) {
		case *Storage_Ref:
			return s.Ref.GetStorage().GetStorageSource()
		default:
			return x.Source
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

type pbMySQLOptions MySQLOptions

func (p *pbMySQLOptions) Reset() {
	(*MySQLOptions)(p).Reset()
}

func (p *pbMySQLOptions) String() string {
	return (*MySQLOptions)(p).String()
}

func (p *pbMySQLOptions) ProtoMessage() {
	(*MySQLOptions)(p).Reset()
}

func (x *MySQLOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewMySQLOptions()
	x.Charset = options.Charset
	x.Collation = options.Collation
	x.MaxIdleConnections = options.MaxIdleConnections
	x.MaxOpenConnections = options.MaxOpenConnections
	x.MaxConnectionLifeTime = options.MaxConnectionLifeTime
	x.TablePrefix = options.TablePrefix
	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbMySQLOptions)(x))
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

//
//type pbGlobalOptions GlobalOptions
//
//func (p *pbGlobalOptions) Reset() {
//	(*GlobalOptions)(p).Reset()
//}
//
//func (p *pbGlobalOptions) String() string {
//	return (*GlobalOptions)(p).String()
//}
//
//func (p *pbGlobalOptions) ProtoMessage() {
//	(*GlobalOptions)(p).Reset()
//}
//
//func (x *GlobalOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
//	options := NewGlobalOptions()
//	x.MaxBodySize = options.MaxBodySize
//	x.MaxUploadSize = options.MaxUploadSize
//	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbGlobalOptions)(x))
//}
//func NewGlobalOptions() *GlobalOptions {
//	opts := &GlobalOptions{
//		MaxUploadSize: &types.UInt32Value{},
//		MaxBodySize:   "5m",
//	}
//	return
//}
