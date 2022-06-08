package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"idas/pkg/logs"
	"idas/pkg/utils/signals"
	"time"
)

func NewMySQLClient(ctx context.Context, options *MySQLOptions) (*Client, error) {
	var m Client
	logger := logs.GetContextLogger(ctx)
	level.Debug(logger).Log("msg", "connect to mysql server",
		"host", options.Host, "username", options.Username,
		"schema", options.Schema,
		"charset", options.Charset,
		"collation", options.Collation)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?parseTime=1&multiStatements=1&charset=%s&collation=%s",
				options.Username,
				options.Password,
				options.Host,
				options.Schema,
				options.Charset,
				options.Collation,
			),
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "t_",
				SingularTable: true,
			},
			Logger: NewLogAdapter(logger),
		},
	)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Errorf("failed to connect to mysql server: [%s@%s]", options.Username, options.Host), "err", err)
		return nil, err
	}

	{
		sqlDB, err := db.DB()
		if err != nil {
			level.Error(logger).Log("msg", fmt.Errorf("failed to connect to mysql server: [%s@%s]", options.Username, options.Host), "err", err)
			return nil, err
		}
		sqlDB.SetMaxIdleConns(int(options.MaxIdleConnections))
		sqlDB.SetConnMaxLifetime(options.GetStdMaxConnectionLifeTime())
		sqlDB.SetMaxOpenConns(int(options.MaxOpenConnections))
	}

	stopCh := signals.SetupSignalHandler(logger)
	stopCh.Add(1)
	go func() {
		<-stopCh.Channel()
		stopCh.WaitRequest()
		if sqlDB, err := db.DB(); err == nil {
			if err = sqlDB.Close(); err != nil {
				level.Warn(logger).Log("msg", fmt.Errorf("failed to close mysql connect: [%s@%s]", options.Username, options.Host), "err", err)
			}
		}
		stopCh.Done()
	}()

	level.Debug(logger).Log("msg", "connected to mysql server",
		"host", options.Host, "username", options.Username,
		"schema", options.Schema,
		"charset", options.Charset,
		"collation", options.Collation)
	m.database = &Database{
		db,
	}
	return &m, nil
}

func (x *MySQLOptions) GetStdMaxConnectionLifeTime() time.Duration {
	if x != nil {
		if duration, err := types.DurationFromProto(x.MaxConnectionLifeTime); err == nil {
			return duration
		}
	}
	return time.Second * 30
}

//
//type pbMySQLOptions MySQLOptions
//
//func (p *pbMySQLOptions) Reset() {
//	(*MySQLOptions)(p).Reset()
//}
//
//func (p *pbMySQLOptions) String() string {
//	return (*MySQLOptions)(p).String()
//}
//
//func (p *pbMySQLOptions) ProtoMessage() {
//	(*MySQLOptions)(p).Reset()
//}
//
//func (m *MySQLOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
//	options := NewMySQLOptions()
//	m.Charset = options.Charset
//	m.Collation = options.Collation
//	m.MaxIdleConnections = options.MaxIdleConnections
//	m.MaxOpenConnections = options.MaxOpenConnections
//	m.MaxConnectionLifeTime = options.MaxConnectionLifeTime
//	m.TablePrefix = options.TablePrefix
//	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbMySQLOptions)(m))
//}
//
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

type MySQLClient struct {
	*Client
	options *MySQLOptions
}

// Merge implement proto.Merger
func (c *MySQLClient) Merge(src proto.Message) {
	if s, ok := src.(*MySQLClient); ok {
		c.options = s.options
		c.Client = s.Client
	}
}

func (c MySQLClient) Options() MySQLOptions {
	return *c.options
}

func (c *MySQLClient) SetOptions(o *MySQLOptions) {
	c.options = o
}

// String implement proto.Message
func (c MySQLClient) String() string {
	return c.options.String()
}

// ProtoMessage implement proto.Message
func (c *MySQLClient) ProtoMessage() {
	c.options.ProtoMessage()
}

// Reset *implement proto.Message*
func (c *MySQLClient) Reset() {
	c.options.Reset()
}

func (c MySQLClient) Marshal() ([]byte, error) {
	return proto.Marshal(c.options)
}

func (c *MySQLClient) Unmarshal(data []byte) (err error) {
	if c.options == nil {
		c.options = NewMySQLOptions()
	}
	if err = proto.Unmarshal(data, c.options); err != nil {
		return err
	}
	if c.Client, err = NewMySQLClient(context.Background(), c.options); err != nil {
		return err
	}
	return
}

func (c MySQLClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.options)
}

func (c *MySQLClient) UnmarshalJSON(data []byte) (err error) {
	if c.options == nil {
		c.options = NewMySQLOptions()
	}
	if err = json.Unmarshal(data, c.options); err != nil {
		return err
	}
	if c.Client, err = NewMySQLClient(context.Background(), c.options); err != nil {
		return err
	}
	return
}
