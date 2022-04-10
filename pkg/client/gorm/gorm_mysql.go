package gorm

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
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
		return nil, fmt.Errorf("连接MySQL数据库[%s@%s]失败(%s)", options.Username, options.Host, err)
	}

	{
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("连接MySQL数据库[%s@%s]失败(%s)", options.Username, options.Host, err)
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
				level.Warn(logger).Log("msg", "关闭MySQL数据库连接失败", "err", err)
			}
		}
		stopCh.Done()
	}()

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
