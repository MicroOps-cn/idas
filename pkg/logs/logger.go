package logs

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-logfmt/logfmt"
	uuid "github.com/satori/go.uuid"

	"idas/pkg/global"
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

func (l *idasLogger) Log(keyvals ...interface{}) error {
	ll := &idasLog{otherKeyMaxLen: 18, caller: log.DefaultCaller}
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
		ll.traceId = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	}
	if ll.level == nil {
		ll.level = LevelInfo
	}
	if ll.caller == nil {
		_, file, line, _ := runtime.Caller(5)
		ll.caller = file + ":" + strconv.Itoa(line)
	}
	if ll.ts == nil {
		ll.ts = timestampFormat()
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
				copy(title[idx:len(ll.title)+idx], []byte(ll.title))
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
func NewIdasLogger(w io.Writer) log.Logger {
	return &idasLogger{w}
}
