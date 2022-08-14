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
	"fmt"

	"github.com/go-kit/log/level"
	"gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
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
