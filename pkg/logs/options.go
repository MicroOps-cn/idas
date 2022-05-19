package logs

import (
	"github.com/go-kit/log"
	"idas/pkg/global"
	"runtime"
	"strings"
)

type Option func(l log.Logger) log.Logger

func WithCaller(layer int) Option {
	return func(l log.Logger) log.Logger {
		return log.With(l, global.CallerName, Caller(layer))
	}
}

func Method(skip ...int) Option {
	pc := make([]uintptr, 1)
	if len(skip) > 0 {
		runtime.Callers(skip[0], pc)
	} else {
		runtime.Callers(2, pc)
	}
	funcName := strings.SplitAfterN(runtime.FuncForPC(pc[0]).Name(), ".", 2)
	return func(l log.Logger) log.Logger {
		return log.With(l, "method", funcName[len(funcName)-1])
	}
}
