package gorm

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"
	"gorm.io/gorm"

	"idas/pkg/global"
	"idas/pkg/logs"
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
	if conn := ctx.Value(global.GormConnName); conn != nil {
		switch db := conn.(type) {
		case *Database:
			return &Database{DB: db.Session(session)}
		case *gorm.DB:
			return &Database{DB: db.Session(session)}
		default:
			level.Warn(logger).Log("msg", "未知的上下文属性(global.GormConnName)值", global.GormConnName, fmt.Sprintf("%#v", conn))
		}
	}
	return &Database{DB: c.database.Session(session).WithContext(ctx)}
}
