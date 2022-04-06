package capacity

import (
	"fmt"
)

type Capacities int64

const (
	Byte     Capacities = 1
	Kilobyte            = Byte << 10
	Megabyte            = Kilobyte << 10
	Gigabyte            = Megabyte << 10
	Terabyte            = Gigabyte << 10
)

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

var magicUnit = []string{"B", "KB", "MB", "GB", "TB"}

func (c Capacities) String() string {
	var buf [32]byte
	w := len(buf)

	u := uint64(c)
	neg := c < 0
	if neg {
		u = -u
	}

	for _, unit := range magicUnit {
		if u > 0 {
			w -= len(unit)
			for idx, r := range unit {
				buf[w+idx] = byte(r)
			}
			if u%1024 > 0 {
				w = fmtInt(buf[:w], u%1024)
			} else {
				w += len(unit)
			}
			u /= 1024
		} else {
			break
		}
		fmt.Println(w, unit)
	}

	if neg {
		w--
		buf[w] = '-'
	}
	return string(buf[w:])
}
