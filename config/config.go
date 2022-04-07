package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"idas/pkg/utils/capacity"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/afero"
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
	if x.Storage == nil || x.Storage.Default == nil {
		return fmt.Errorf("default storage is null")
	}
	if len(x.Storage.User) == 0 {
		x.Storage.User = append(x.Storage.User, x.Storage.Default)
	}
	if len(x.Storage.App) == 0 {
		x.Storage.App = append(x.Storage.App, x.Storage.Default)
	}
	if x.Storage.Session == nil {
		x.Storage.Session = x.Storage.Default
	}
	var storages = append(append(x.Storage.User, x.Storage.App...), x.Storage.Session, x.Storage.Default)
	for _, storage := range storages {
		switch s := storage.Source.(type) {
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

type pbGlobalOptions GlobalOptions

func (p *pbGlobalOptions) Reset() {
	(*GlobalOptions)(p).Reset()
}

func (p *pbGlobalOptions) String() string {
	return (*GlobalOptions)(p).String()
}

func (p *pbGlobalOptions) ProtoMessage() {
	(*GlobalOptions)(p).Reset()
}

func (x *GlobalOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewGlobalOptions()
	x.MaxBodySize = options.MaxBodySize
	x.MaxUploadSize = options.MaxUploadSize
	x.UploadPath = "uploads"
	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbGlobalOptions)(x))
}

const defaultMaxUploadSize = 1 << 20 * 10
const defaultMaxBodySize = 1 << 20 * 5

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		MaxUploadSize: capacity.NewCapacity(defaultMaxUploadSize),
		MaxBodySize:   capacity.NewCapacity(defaultMaxBodySize),
	}
}

func (x *Config) GetUploadDir() (afero.Fs, error) {
	var (
		ws         afero.Fs
		uploadPath = "uploads"
	)
	if x.Global != nil {
		if len(x.Global.UploadPath) > 0 {
			uploadPath = x.Global.UploadPath
			if x.Global.UploadPath[0] == '/' {
				ws = afero.NewOsFs()
			}
		}
	}
	if ws == nil {
		ws = x.GetWorkspace()
	}
	if stat, err := ws.Stat(uploadPath); os.IsNotExist(err) {
		if err = os.MkdirAll(uploadPath, 0755); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("path `%s` is not directory", x.Global.UploadPath)
	}
	return afero.NewBasePathFs(ws, uploadPath), nil
}
func (x *Config) GetWorkspace() afero.Fs {
	if x.Global != nil {
		if len(x.Global.Workspace) != 0 {
			return afero.NewBasePathFs(afero.NewOsFs(), x.Global.Workspace)
		}
	}
	return nil
}

func (x *Config) SetWorkspace(path string) {
	if x.Global == nil {
		x.Global = new(GlobalOptions)
	}
	x.Global.Workspace = path
}

//func (x *Config) GetMaxBodySize(path string) {
//	if x.Global == nil {
//		x.Global = new(GlobalOptions)
//	}
//	x.Global.Workspace = path
//}
