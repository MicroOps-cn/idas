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
	"time"

	log2 "github.com/MicroOps-cn/fuck/log"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/MicroOps-cn/idas/pkg/logs"
)

type logContext struct {
	logger        kitlog.Logger
	slowThreshold time.Duration
}

func (l *logContext) LogMode(lvl logger.LogLevel) logger.Interface {
	var filter kitlog.Logger
	switch lvl {
	case logger.Silent:
		filter = level.NewFilter(l.logger, level.AllowNone())
	case logger.Info:
		filter = level.NewFilter(l.logger, level.AllowInfo())
	case logger.Warn:
		filter = level.NewFilter(l.logger, level.AllowWarn())
	case logger.Error:
		filter = level.NewFilter(l.logger, level.AllowError())
	default:
		filter = l.logger
	}

	return NewLogAdapter(filter, l.slowThreshold)
}

func (l logContext) Info(_ context.Context, msg string, data ...interface{}) {
	level.Info(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Warn(_ context.Context, msg string, data ...interface{}) {
	level.Warn(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Error(_ context.Context, msg string, data ...interface{}) {
	level.Error(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && err != gorm.ErrRecordNotFound:
		sql, rows := fc()
		level.Error(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", "SQL execution exception", logs.WrapKeyName("errorMsg"), err, logs.WrapKeyName("sql"), sql, logs.WrapKeyName("execTime"), float64(elapsed.Nanoseconds())/1e6, logs.WrapKeyName("rowReturnCount"), rows)
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		sql, rows := fc()
		level.Warn(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", "exec SQL query", logs.WrapKeyName("sql"), sql, logs.WrapKeyName("execTime"), float64(elapsed.Nanoseconds())/1e6, logs.WrapKeyName("rowReturnCount"), rows)
	default:
		sql, rows := fc()
		level.Debug(l.logger).Log(log2.CallerName, logs.Relative(utils.FileWithLineNum()), "msg", "exec SQL query", logs.WrapKeyName("sql"), sql, logs.WrapKeyName("execTime"), float64(elapsed.Nanoseconds())/1e6, logs.WrapKeyName("rowReturnCount"), rows)
	}
}

func NewLogAdapter(l kitlog.Logger, slowThreshold time.Duration) logger.Interface {
	return &logContext{logger: l, slowThreshold: slowThreshold}
}

var _ logger.Interface = new(logContext)
