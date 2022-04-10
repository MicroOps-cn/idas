package gorm

import (
	"bytes"
	"context"
	"fmt"
	"gorm.io/gorm/schema"

	"github.com/go-kit/log/level"
	"github.com/golang/protobuf/jsonpb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func (x *SQLiteOptions) UnmarshalJSONPB(unmarshaller *jsonpb.Unmarshaler, b []byte) error {
	options := NewSQLiteOptions()
	x.Path = options.Path
	return unmarshaller.Unmarshal(bytes.NewReader(b), (*pbSQLiteOptions)(x))
}

func NewSQLiteOptions() *SQLiteOptions {
	return &SQLiteOptions{
		Path: "idas.db",
	}
}
