package logs

import (
	"github.com/go-kit/log"
	"idas/pkg/global"
)

type Option func(l log.Logger) log.Logger

func Caller(layer int) Option {
	return func(l log.Logger) log.Logger {
		return log.With(l, global.CallerName, log.Caller(layer))
	}
}
