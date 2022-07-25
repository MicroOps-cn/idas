package gorm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"idas/pkg/logs"
	"idas/pkg/utils/signals"
)

func NewSQLiteClient(ctx context.Context, options *SQLiteOptions) (*Client, error) {
	var m Client
	logger := logs.GetContextLogger(ctx)
	db, err := gorm.Open(sqlite.Open(options.Path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: NewLogAdapter(logger),
	})
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库[%s]失败: %s", options.Path, err)
	}

	stopCh := signals.SetupSignalHandler(logger)
	stopCh.Add(1)
	go func() {
		<-stopCh.Channel()
		stopCh.WaitRequest()
		if sqlDB, err := db.DB(); err == nil {
			if err = sqlDB.Close(); err != nil {
				level.Warn(logger).Log("msg", "关闭SQLite数据库连接失败", "err", err)
			}
		}
		stopCh.Done()
	}()

	m.database = &Database{DB: db}
	return &m, nil
}

type pbSQLiteOptions SQLiteOptions

func (p *pbSQLiteOptions) Reset() {
	(*SQLiteOptions)(p).Reset()
}

func (p *pbSQLiteOptions) String() string {
	return (*SQLiteOptions)(p).String()
}

func (p *pbSQLiteOptions) ProtoMessage() {
	(*SQLiteOptions)(p).Reset()
}

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
