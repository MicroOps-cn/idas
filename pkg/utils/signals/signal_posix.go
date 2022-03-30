//go:build !windows

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
