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

package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/MicroOps-cn/fuck/clients/gorm"
	oauth2 "github.com/MicroOps-cn/idas/pkg/client/oauth2"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/utils/capacity"
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
	var curPath, downPath string
	if idx := strings.IndexAny(path, ".["); idx >= 0 {
		downPath = path[idx:]
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
		curPath = path[:idx]
	} else {
		curPath = path
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
		if jsonName == curPath {
			return ref(downPath, sv)
		}
	}
	return nil
}

func (x *Storage) findRef(path string, root interface{}) error {
	target := ref(path, reflect.ValueOf(root))
	var tmpStorage *Storage
	switch s := target.(type) {
	case Storage:
		if s.GetRef() != nil {
			return x.findRef(s.GetRef().Path, root)
		}
		tmpStorage = proto.Clone(&s).(*Storage)
	case *Storage:
		if s.GetRef() != nil {
			return x.findRef(s.GetRef().Path, root)
		}
		tmpStorage = proto.Clone(s).(*Storage)
	default:
		return fmt.Errorf("unknown ref: %s(%T)", path, target)
	}
	x.Source = tmpStorage.GetStorageSource()
	return nil
}

func (x *Config) Init(_ log.Logger) error {
	if x.Storage == nil {
		x.Storage = &Storages{}
	}
	if x.Storage.Default == nil {
		o := gorm.NewSQLiteOptions()
		x.Storage.Default = &Storage{
			Name: "default",
			Source: &Storage_Sqlite{
				Sqlite: w.M(gorm.NewSQLiteClient(context.Background(), o)),
			},
		}
	}
	if x.Storage.User == nil {
		x.Storage.User = x.Storage.Default
	}
	if x.Storage.Session == nil {
		x.Storage.Session = x.Storage.Default
	}
	if x.Storage.Logging == nil {
		x.Storage.Logging = x.Storage.Default
	}
	storages := []*Storage{x.Storage.User, x.Storage.Session, x.Storage.Default}
	for _, storage := range storages {
		switch s := storage.Source.(type) {
		case *Storage_Ref:
			if s.Ref.Storage == nil {
				s.Ref.Storage = new(Storage)
			}
			err := storage.findRef(s.Ref.Path, x)
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

type pbGlobalOptions GlobalOptions

func (p *pbGlobalOptions) Reset() {
	(*GlobalOptions)(p).Reset()
}

func (p *pbGlobalOptions) String() string {
	return (*GlobalOptions)(p).String()
}

func (p *pbGlobalOptions) ProtoMessage() {
	(*GlobalOptions)(p).ProtoMessage()
}

func (x *GlobalOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewGlobalOptions()
	x.MaxBodySize = options.MaxBodySize
	x.MaxUploadSize = options.MaxUploadSize
	x.JwtSecret = options.JwtSecret
	x.AppName = options.AppName
	x.UploadPath = options.UploadPath
	x.Title = options.Title
	err := unmarshaller.Unmarshal(bytes.NewReader(b), (*pbGlobalOptions)(x))
	if err != nil {
		return err
	}
	if len(x.JwtSecret) == 0 {
		return fmt.Errorf("`global.jwt_secret` cannot be empty")
	}
	if len(x.Secret) == 0 {
		return fmt.Errorf("`global.secret` cannot be empty")
	}
	return nil
}

const (
	defaultMaxUploadSize = 1 << 20 * 10
	defaultMaxBodySize   = 1 << 20 * 5
)

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		MaxUploadSize: capacity.NewCapacity(defaultMaxUploadSize),
		MaxBodySize:   capacity.NewCapacity(defaultMaxBodySize),
		JwtSecret:     rand.String(128),
		AppName:       global.AppName,
		Title:         "IDAS",
		UploadPath:    "uploads",
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
		//nolint:gofumpt
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

func (x *Config) GetAppName() string {
	if x.Global != nil {
		if len(x.Global.AppName) != 0 {
			return x.Global.AppName
		}
	}
	return "idas"
}

func (x *Config) GetOAuthOptions(id string) *oauth2.Options {
	for _, option := range x.GetGlobal().GetOauth2() {
		if option.Id == id {
			return option
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

func (c *RuntimeConfig) GetPasswordFailedLockConfig() (failedSec, failedThreshold int64) {
	sec := c.GetSecurity()
	if sec != nil {
		failedThreshold = int64(sec.GetPasswordFailedLockThreshold())
		failedSec = int64(sec.GetPasswordFailedLockDuration()) * 60
		return failedSec, failedThreshold
	}
	return 0, 0
}

func (c *RuntimeConfig) GetLoginSessionInactivityTime() uint32 {
	sec := c.GetSecurity()
	inactiveTime := sec.GetLoginSessionInactivityTime()
	maxTime := sec.GetLoginSessionMaxTime()
	if inactiveTime > maxTime {
		return maxTime
	}
	if inactiveTime > 0 {
		return inactiveTime
	}
	return 30 * 24 // 30天
}

func (c *RuntimeConfig) GetLoginSessionMaxTime() uint32 {
	maxTime := c.GetSecurity().GetLoginSessionMaxTime()
	if maxTime > 0 {
		return maxTime
	}
	return 30 * 24 // 30天
}
