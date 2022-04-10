package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type logContext struct {
	logger        log.Logger
	SlowThreshold time.Duration
}

func (l *logContext) LogMode(lvl logger.LogLevel) logger.Interface {
	var filter log.Logger
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

	return NewLogAdapter(filter)
}

func (l logContext) Info(_ context.Context, msg string, data ...interface{}) {
	level.Info(l.logger).Log("caller", utils.FileWithLineNum(), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Warn(_ context.Context, msg string, data ...interface{}) {
	level.Warn(l.logger).Log("caller", utils.FileWithLineNum(), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Error(_ context.Context, msg string, data ...interface{}) {
	level.Error(l.logger).Log("caller", utils.FileWithLineNum(), "msg", fmt.Sprintf(msg, data...))
}

func (l logContext) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && err != gorm.ErrRecordNotFound:
		sql, rows := fc()
		level.Error(l.logger).Log("caller", utils.FileWithLineNum(), "msg", "SQL execution exception", "[ErrorMsg]", err, "[sql]", sql, "[ExecTime]", float64(elapsed.Nanoseconds())/1e6, "[RowReturnCount]", rows)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		level.Warn(l.logger).Log("caller", utils.FileWithLineNum(), "msg", "exec SQL query", "[sql]", sql, "[ExecTime]", float64(elapsed.Nanoseconds())/1e6, "[RowReturnCount]", rows)
	default:
		sql, rows := fc()
		level.Debug(l.logger).Log("caller", utils.FileWithLineNum(), "msg", "exec SQL query", "[sql]", sql, "[ExecTime]", float64(elapsed.Nanoseconds())/1e6, "[RowReturnCount]", rows)
	}
}

func NewLogAdapter(l log.Logger) logger.Interface {
	return &logContext{logger: l}
}

var _ logger.Interface = new(logContext)
