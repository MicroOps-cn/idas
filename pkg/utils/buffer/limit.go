package buffer

import "io"

type LimitWriterOption func(*LimitedWriter)

// LimitWriter returns a Writer that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedWriter.
func LimitWriter(w io.Writer, n int64, options ...LimitWriterOption) io.Writer {
	writer := &LimitedWriter{W: w, N: n, MaxN: n}
	for _, option := range options {
		option(writer)
	}
	return writer
}

func LimitWriterIgnoreError(lw *LimitedWriter) {
	lw.ignoreError = true
}

// A LimitedWriter reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type LimitedWriter struct {
	W           io.Writer // underlying reader
	N           int64     // max bytes remaining
	MaxN        int64     // max bytes remaining
	ignoreError bool
}

func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	var c = len(p)
	if l.N <= 0 {
		err = io.EOF
	} else {
		if int64(len(p)) > l.N {
			p = p[0:l.N]
		}
		n, err = l.W.Write(p)
		l.N -= int64(n)
	}
	if l.ignoreError {
		return c, nil
	}

	return
}
