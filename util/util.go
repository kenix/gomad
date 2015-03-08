package util

import (
	"io"
	"log"
	"math/rand"
	"time"
)

func Shuffle(a []int64) {
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func GoNonBlocking(f func() interface{}) (ch chan interface{}) {
	go func() {
		select {
		case ch <- f():
		default:
		}
	}()
	return
}

var Lt *log.Logger
var Li *log.Logger
var Lw *log.Logger
var Le *log.Logger
var Lf *log.Logger
var Lc *log.Logger

func InitLoggers(w io.Writer) {
	Lt = log.New(w, "T ", log.LstdFlags|log.Lshortfile)
	Li = log.New(w, "I ", log.LstdFlags|log.Lshortfile)
	Lw = log.New(w, "W ", log.LstdFlags|log.Lshortfile)
	Le = log.New(w, "E ", log.LstdFlags|log.Lshortfile)
	Lf = log.New(w, "F ", log.LstdFlags|log.Lshortfile)
	Lc = log.New(w, "C ", log.LstdFlags|log.Lshortfile)
}

type Cue struct{}

func LeadingZeros(d interface{}) int {
	switch t := d.(type) {
	case uint8:
		return lzs(uint64(t), 8)
	case int8:
		return lzs(uint64(t), 8)
	case int16:
		return lzs(uint64(t), 16)
	case uint16:
		return lzs(uint64(t), 16)
	case int32:
		return lzs(uint64(t), 32)
	case uint32:
		return lzs(uint64(t), 32)
	case int64:
		return lzs(uint64(t), 64)
	case uint64:
		return lzs(uint64(t), 64)
	}
	return -1 // unreachable code
}

func lzs(x uint64, s uint) int {
	u := uint64(1 << (s - 1))
	for i := uint(0); i < s; i++ {
		if (u>>i)&x > 0 {
			return int(i)
		}
	}
	return int(s)
}

func Size(d interface{}) int {
	switch d.(type) {
	case uint8, int8: // also applied to byte
		return 8
	case uint16, int16:
		return 16
	case int32, uint32: // also applied to rune
		return 32
	case int64, uint64:
		return 64
	default:
		return -1
	}
}
