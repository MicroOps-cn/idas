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

package gorm

import (
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MicroOps-cn/fuck/log"
	gosqlite "github.com/glebarez/go-sqlite"
	"github.com/glebarez/sqlite"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/MicroOps-cn/idas/api"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

func init() {
	gosqlite.MustRegisterDeterministicScalarFunction("from_base64", 1, func(ctx *gosqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		switch argTyped := args[0].(type) {
		case string:
			return base64.StdEncoding.DecodeString(argTyped)
		default:
			return nil, fmt.Errorf("unsupported type: %T", args[0])
		}
	})
}

func NewSQLiteClient(ctx context.Context, options *SQLiteOptions) (clt *Client, err error) {
	clt = new(Client)
	logger := log.GetContextLogger(ctx)
	if options.SlowThreshold != nil {
		clt.slowThreshold, err = types.DurationFromProto(options.SlowThreshold)
		if err != nil {
			level.Error(logger).Log("msg", fmt.Sprintf("failed to connect to SQLite database: %s", options.Path), "err", fmt.Errorf("`slow_threshold` option is invalid: %s", err))
			return nil, err
		}
	}

	level.Debug(logger).Log("msg", "connect to sqlite", "dsn", options.Path)
	db, err := gorm.Open(sqlite.Open(options.Path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: NewLogAdapter(logger, clt.slowThreshold),
	})
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("failed to connect to SQLite database: %s", options.Path), "err", err)
		return nil, errors.WithServerError(http.StatusInternalServerError, err, fmt.Sprintf("failed to connect to SQLite database: %s", options.Path))
	}

	stopCh := signals.SetupSignalHandler(logger)
	stopCh.Add(1)
	go func() {
		<-stopCh.Channel()
		stopCh.WaitRequest()
		if sqlDB, err := db.DB(); err == nil {
			if err = sqlDB.Close(); err != nil {
				level.Warn(logger).Log("msg", "Failed to close SQLite database", "err", err)
			}
		}
		level.Debug(logger).Log("msg", "Sqlite connect closed")
		stopCh.Done()
	}()

	clt.database = &Database{DB: db}
	return clt, nil
}

//
//type pbSQLiteOptions SQLiteOptions
//
//func (p *pbSQLiteOptions) Reset() {
//	(*SQLiteOptions)(p).Reset()
//}
//
//func (p *pbSQLiteOptions) String() string {
//	return (*SQLiteOptions)(p).String()
//}
//
//func (p *pbSQLiteOptions) ProtoMessage() {
//	(*SQLiteOptions)(p).Reset()
//}

func NewSQLiteOptions() *SQLiteOptions {
	return &SQLiteOptions{
		Path:        "idas.db",
		TablePrefix: "t_",
	}
}

type SQLiteClient struct {
	*Client
	options *SQLiteOptions
}

var _ api.CustomType = &SQLiteClient{}

// Merge implement proto.Merger
func (c *SQLiteClient) Merge(src proto.Message) {
	if s, ok := src.(*SQLiteClient); ok {
		c.options = s.options
		c.Client = s.Client
	}
}

// Reset *implement proto.Message*
func (c *SQLiteClient) Reset() {
	c.options.Reset()
}

// String implement proto.Message
func (c SQLiteClient) String() string {
	return c.options.String()
}

// ProtoMessage implement proto.Message
func (c *SQLiteClient) ProtoMessage() {
	c.options.ProtoMessage()
}

func (c SQLiteClient) Marshal() ([]byte, error) {
	return proto.Marshal(c.options)
}

func (c *SQLiteClient) Unmarshal(data []byte) (err error) {
	if c.options == nil {
		c.options = NewSQLiteOptions()
	}
	if err = proto.Unmarshal(data, c.options); err != nil {
		return err
	}
	if c.Client, err = NewSQLiteClient(context.Background(), c.options); err != nil {
		return err
	}
	return
}

var _ proto.Unmarshaler = &SQLiteClient{}

func (c SQLiteClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.options)
}

func (c *SQLiteClient) UnmarshalJSON(data []byte) (err error) {
	if c.options == nil {
		c.options = NewSQLiteOptions()
	}
	if err = json.Unmarshal(data, c.options); err != nil {
		return err
	}
	if c.Client, err = NewSQLiteClient(context.Background(), c.options); err != nil {
		return err
	}
	return
}
