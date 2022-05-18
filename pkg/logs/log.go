// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logs defines standardised ways to initialize Go kit loggers
// across Prometheus components.
// It should typically only ever be imported by main packages.
package logs

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"idas/pkg/global"
)

// This timestamp format differs from RFC3339Nano by using .000 instead
// of .999999999 which changes the timestamp from 9 variable to 3 fixed
// decimals (.130 instead of .130987456).
var timestampFormat = log.TimestampFormat(
	func() time.Time { return time.Now().UTC() },
	"2006-01-02T15:04:05.000Z07:00",
)

// AllowedLevel is a settable identifier for the minimum level a log entry
// must be have.
type AllowedLevel string

func (l AllowedLevel) getOption() level.Option {
	switch l {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	default:
		return level.AllowWarn()
	}
}

func (l *AllowedLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	type plain string
	if err := unmarshal((*plain)(&s)); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	var lo AllowedLevel
	if err := lo.Set(s); err != nil {
		return err
	}
	*l = lo
	return nil
}

func (l AllowedLevel) String() string {
	return string(l)
}

func (l AllowedLevel) Valid() error {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return nil
	default:
		return errors.Errorf("unrecognized log level %s", l)
	}
}

// Set updates the value of the allowed level.
func (l *AllowedLevel) Set(s string) error {
	lvl := AllowedLevel(s)
	if err := lvl.Valid(); err != nil {
		return err
	}
	*l = AllowedLevel(s)
	return nil
}

const (
	LevelDebug AllowedLevel = "debug"
	LevelInfo  AllowedLevel = "info"
	LevelWarn  AllowedLevel = "warn"
	LevelError AllowedLevel = "error"
)

// AllowedFormat is a settable identifier for the output format that the logger can have.
type AllowedFormat string

func (f AllowedFormat) String() string {
	return string(f)
}

func (f AllowedFormat) Valid() error {
	switch f {
	case FormatIdas, FormatLogfmt, FormatJSON:
		return nil
	default:
		return errors.Errorf("unrecognized log format %s", f)
	}
}

// Set updates the value of the allowed format.
func (f *AllowedFormat) Set(s string) error {
	format := AllowedFormat(s)
	if err := format.Valid(); err != nil {
		return err
	}
	*f = format
	return nil
}

const (
	FormatIdas   AllowedFormat = "idas"
	FormatJSON   AllowedFormat = "json"
	FormatLogfmt AllowedFormat = "logfmt"
)

// Config is a struct containing configurable settings for the logger
type Config struct {
	Level  *AllowedLevel
	Format *AllowedFormat
}

func MustNewConfig(level string, format string) *Config {
	cfg := &Config{Level: new(AllowedLevel), Format: new(AllowedFormat)}
	if err := cfg.Level.Set(level); err != nil {
		panic(err)
	}
	if err := cfg.Format.Set(format); err != nil {
		panic(err)
	}
	return cfg
}

// New returns a new leveled oklog logger. Each logged line will be annotated
// with a timestamp. The output always goes to stderr.
func New(config *Config) log.Logger {
	var l log.Logger
	switch *config.Format {
	case FormatLogfmt:
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	case FormatJSON:
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	default:
		l = NewIdasLogger(log.NewSyncWriter(os.Stderr))
	}

	if config.Level != nil {
		l = level.NewFilter(l, config.Level.getOption())
	}
	l = log.With(l, "ts", timestampFormat, global.CallerName, log.DefaultCaller)
	return l
}

// NewDynamic returns a new leveled logger. Each logged line will be annotated
// with a timestamp. The output always goes to stderr. Some properties can be
// changed, like the level.
func newDynamic(config *Config) *logger {
	var l log.Logger
	if config.Format != nil && *config.Format == "json" {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}
	l = log.With(l, "ts", timestampFormat, global.CallerName, log.DefaultCaller)

	lo := &logger{
		base:    l,
		leveled: l,
	}
	if config.Level != nil {
		lo.SetLevel(config.Level)
	}
	return lo
}

type logger struct {
	base         log.Logger
	leveled      log.Logger
	currentLevel *AllowedLevel
	mtx          sync.Mutex
}

// Log implements logger.Log.
func (l *logger) Log(keyvals ...interface{}) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.leveled.Log(keyvals...)
}

// SetLevel changes the log level.
func (l *logger) SetLevel(lvl *AllowedLevel) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if lvl != nil {
		if l.currentLevel != nil && *l.currentLevel != *lvl {
			_ = l.base.Log("msg", "Log level changed", "prev", l.currentLevel, "current", lvl)
		}
		l.currentLevel = lvl
	}
	l.leveled = level.NewFilter(l.base, lvl.getOption())
}

var rootLogger log.Logger

func SetRootLogger(logger log.Logger) {
	rootLogger = logger
}

func GetRootLogger() log.Logger {
	if rootLogger == nil {
		panic("root logger is uninitialized")
	}
	return rootLogger
}

func NewTraceLogger() log.Logger {
	traceId := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	return log.With(rootLogger, global.TraceIdName, traceId)
}

func GetContextLogger(ctx context.Context, options ...Option) log.Logger {
	l, ok := ctx.Value(global.LoggerName).(log.Logger)
	if !ok {
		l = NewTraceLogger()
	}
	for _, option := range options {
		l = option(l)
	}
	return log.With(l)
}

type WriterAdapter struct {
	l               log.Logger
	msgKey          string
	prefix          string
	joinPrefixToMsg bool
}

func (a WriterAdapter) Write(p []byte) (n int, err error) {
	a.l.Log(a.msgKey, a.handleMessagePrefix(string(p)))
	return len(p), nil
}

func (a WriterAdapter) handleMessagePrefix(msg string) string {
	if a.prefix == "" {
		return msg
	}

	msg = strings.TrimPrefix(msg, a.prefix)
	if a.joinPrefixToMsg {
		msg = a.prefix + msg
	}
	return msg
}

func MessageKey(key string) WriterAdapterOption {
	return func(a *WriterAdapter) { a.msgKey = key }
}

func Prefix(prefix string, joinPrefixToMsg bool) WriterAdapterOption {
	return func(a *WriterAdapter) { a.prefix = prefix; a.joinPrefixToMsg = joinPrefixToMsg }
}

type WriterAdapterOption func(*WriterAdapter)

func NewWriterAdapter(logger log.Logger, options ...WriterAdapterOption) io.Writer {
	adapter := &WriterAdapter{l: logger, msgKey: "msg"}
	for _, option := range options {
		option(adapter)
	}
	return adapter
}
