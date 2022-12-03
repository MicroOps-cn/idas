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
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/MicroOps-cn/fuck/log"
	kitlog "github.com/go-kit/log"
	"github.com/go-logfmt/logfmt"

	"github.com/MicroOps-cn/idas/pkg/global"
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

var sourceDir = func() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return strings.TrimSuffix(currentFile, "pkg/logs/logger.go")
}()

func Caller(depth int) kitlog.Valuer {
	return func() interface{} {
		//for i := 1; i < 15; i++ {
		//	_, f, l, _ := runtime.Caller(i)
		//	if strings.HasPrefix(f, sourceDir) {
		//		fmt.Printf("%d/%d => %s:%d\n", depth, i, f, l)
		//	}
		//}
		_, file, line, _ := runtime.Caller(depth)
		return strings.TrimPrefix(file, sourceDir) + ":" + strconv.Itoa(line)
	}
}

var DefaultCaller = Caller(3)

func (l *idasLogger) Log(keyvals ...interface{}) error {
	ll := &idasLog{otherKeyMaxLen: 18, caller: DefaultCaller}
	for i := 0; ; {
		v := keyvals[i+1]
		if k, ok := keyvals[i].(string); ok {
			if k == "level" {
				ll.level = v
			} else if k == "ts" {
				ll.ts = v
			} else if k == global.CallerName {
				ll.caller = v
			} else if k == "msg" {
				ll.msg = v
			} else if k == "title" {
				ll.title = fmt.Sprintf("%s", v)
			} else if k == global.TraceIdName {
				ll.traceId = v
			} else {
				if len(k) > 0 && k[0] == '[' && k[len(k)-1] == ']' {
					ll.other = append(ll.other, logKvPair{key: k, val: v})
					if len(k) > ll.otherKeyMaxLen {
						ll.otherKeyMaxLen = len(k)
					}
				} else {
					ll.kvs = append(ll.kvs, k, v)
				}
			}
		} else {
			ll.kvs = append(ll.kvs, k, v)
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
		fmt.Printf("格式化异常==>%s(%s)\n", buffer.String(), err)
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

const FormatIDAS log.AllowedFormat = "idas"

func init() {
	log.RegisterLogFormat(FormatIDAS, NewIdasLogger)
}

func Relative(file string) string {
	return strings.TrimPrefix(file, sourceDir)
}
