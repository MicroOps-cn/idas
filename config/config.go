package config

import (
	"encoding/json"
	"time"

	"github.com/go-kit/log"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
)

func (x *Config) Init(logger log.Logger) error {
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
