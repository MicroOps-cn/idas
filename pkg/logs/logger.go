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

package logs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/MicroOps-cn/fuck/log"
	kitlog "github.com/go-kit/log"
	"github.com/go-logfmt/logfmt"
)

const titleBg = "---------------------------------------------------------------------------------------\n"

type logfmtEncoder struct {
	*logfmt.Encoder
	buf bytes.Buffer
}

func (l *logfmtEncoder) Reset() {
	l.Encoder.Reset()
	l.buf.Reset()
}

var idasEncoderPool = sync.Pool{
	New: func() interface{} {
		var enc logfmtEncoder
		enc.Encoder = logfmt.NewEncoder(&enc.buf)
		return &enc
	},
}

type titleKey struct{}

func (titleKey) String() string {
	return "title"
}

var TitleKey = &titleKey{}

type wrapKeyName string

func (n wrapKeyName) String() string {
	return string(n)
}

func WrapKeyName(name string) fmt.Stringer {
	return wrapKeyName(name)
}

type logKvPair struct {
	key string
	val interface{}
}

type idasLog struct {
	level          interface{}
	ts             interface{}
	caller         interface{}
	traceId        interface{}
	msg            interface{}
	title          string
	kvs            []interface{}
	other          []logKvPair
	otherKeyMaxLen int
}

type idasLogger struct {
	w io.Writer
}

func (l *idasLogger) encodeKeyvals(keyvals ...interface{}) ([]byte, error) {
	enc := idasEncoderPool.Get().(*logfmtEncoder)
	enc.Reset()
	defer idasEncoderPool.Put(enc)

	if err := enc.EncodeKeyvals(keyvals...); err != nil {
		return nil, err
	}

	// Add newline to the end of the buffer
	if err := enc.EndRecord(); err != nil {
		return nil, err
	}
	return enc.buf.Bytes(), nil
}

func (l *idasLogger) Log(keyvals ...interface{}) error {
	ll := &idasLog{otherKeyMaxLen: 18, caller: log.DefaultCaller}
	for i := 0; ; {
		v := keyvals[i+1]
		if keyvals[i] == TitleKey {
			ll.title = fmt.Sprintf("%s", v)
		}
		switch k := keyvals[i].(type) {
		case titleKey, *titleKey:
			ll.title = fmt.Sprintf("%s", v)
		case wrapKeyName, *wrapKeyName:
			key := fmt.Sprintf("[%s]", k)
			ll.other = append(ll.other, logKvPair{key: key, val: v})
			if len(key) > ll.otherKeyMaxLen {
				ll.otherKeyMaxLen = len(key)
			}
		case string:
			if k == "level" {
				ll.level = v
			} else if k == "ts" {
				ll.ts = v
			} else if k == "msg" {
				ll.msg = v
			} else if k == log.TraceIdName {
				ll.traceId = v
			} else {
				ll.kvs = append(ll.kvs, k, v)
			}
		default:
			if k == log.CallerName {
				ll.caller = v
			} else {
				ll.kvs = append(ll.kvs, k, v)
			}
		}
		i += 2
		if i >= len(keyvals) {
			break
		}
	}
	if ll.traceId == nil {
		ll.traceId = log.NewTraceId()
	}
	if ll.level == nil {
		ll.level = log.LevelInfo
	}
	if ll.caller == nil {
		_, file, line, _ := runtime.Caller(5)
		ll.caller = file + ":" + strconv.Itoa(line)
	}
	if ll.ts == nil {
		ll.ts = log.TimestampFormat()
	}
	if ll.msg == nil {
		ll.msg = ""
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("%s [%s] %s %s - %v - ", ll.ts, ll.level, ll.traceId, ll.caller, ll.msg))

	if data, err := l.encodeKeyvals(ll.kvs...); err != nil {
		return err
	} else if _, err = buffer.Write(data); err != nil {
		return err
	} else if len(ll.title) > 0 || len(ll.other) > 0 {
		if len(ll.title) > 0 {
			if len(ll.title) > len(titleBg) {
				buffer.WriteString(ll.title)
			} else {
				title := []byte(titleBg)
				idx := (len(title) - len(ll.title)) / 2
				copy(title[idx:len(ll.title)+idx], ll.title)
				buffer.Write(title)
			}
		}
		for _, v := range ll.other {
			buffer.WriteString(fmt.Sprintf("%-"+strconv.Itoa(ll.otherKeyMaxLen)+"s%v\n", fmt.Sprintf("%s:", v.key), v.val))
		}
		if len(ll.title) > 0 {
			buffer.WriteString(titleBg)
		}
	}
	if _, err := l.w.Write(buffer.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: log=%s,err=%s\n", buffer.String(), err)
	}
	return nil
}

// NewIdasLogger returns a logger that encodes keyvals to the Writer in
// logfmt format. Each log event produces no more than one call to w.Write.
// The passed Writer must be safe for concurrent use by multiple goroutines if
// the returned Logger will be used concurrently.
func NewIdasLogger(w io.Writer) kitlog.Logger {
	return &idasLogger{w}
}

var sourceDir = log.GetSourceCodeDir("pkg/logs/logger.go")

const FormatIDAS log.AllowedFormat = "idas"

func init() {
	log.RegisterLogFormat(FormatIDAS, NewIdasLogger)
	log.SetSourceCodeDir(sourceDir)
}

func Relative(file string) string {
	return strings.TrimPrefix(file, sourceDir)
}
