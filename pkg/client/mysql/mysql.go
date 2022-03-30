package mysql

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"idas/config"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/utils/signals"
)

type Database struct {
	*gorm.DB
}

type Client struct {
	database *Database
}

func (c *Client) Session(ctx context.Context) *Database {
	logger := logs.GetContextLogger(ctx)
	session := &gorm.Session{Logger: NewLogAdapter(logger)}
	if conn := ctx.Value(global.MySQLConnName); conn != nil {
		switch db := conn.(type) {
		case *Database:
			return &Database{DB: db.Session(session)}
		case *gorm.DB:
			return &Database{DB: db.Session(session)}
		default:
			level.Warn(logger).Log("msg", "未知的上下文属性(global.MySQLConnName)值", global.MySQLConnName, fmt.Sprintf("%#v", conn))
		}
	}
	return &Database{DB: c.database.Session(session)}
}

func NewMySQLClient(ctx context.Context, options *config.MySQLOptions) (*Client, error) {
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
