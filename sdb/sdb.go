package sdb

import (
	"errors"
	"io"
	"math"
	"os"
	"sort"

	bb "github.com/kenix/gomad/bytebuffer"
	u "github.com/kenix/gomad/util"
)

type Writer interface {
	io.Closer
	Put(key string, dat []byte) (n int, err error)
}

const (
	MaxKeyLength  = math.MaxUint8 // 255
	MaxDataLength = 1 << 24       // 16M
	size_buf_dat  = MaxDataLength // 16M
	size_page_key = 1 << 12       // 4K
)

type entry struct {
	key    string
	offset int64 // 40 bits used, max 1T
	length int32 // 24 bits used, max 16M, offset, length together 8 bytes
}

type entries []entry

func (es entries) Len() int {
	return len(es)
}

func (es entries) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (es entries) Less(i, j int) bool {
	return entries[i] < entries[j]
}

type wImpl struct {
	f    *os.File
	buf  bb.ByteBuffer
	keys entries
	cur  int64
}

func NewWriter(fn string) (Writer, error) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &wImpl{f, bb.New(size_buf_dat), make([]entry, 0, 0), 0}, nil
}

func (w *wImpl) Close() error {
	defer w.f.Close()
	if err := fwc(w.buf, w.f); err != nil {
		return err
	}
	bls, err := w.persistKeys() // B*-tree level 1 leaves start position
	if err != nil {
		return err
	}
	w.buf.PutUint64(uint64(bls))
	if err := fwc(w.buf, w.f); err != nil {
		return err
	}
}

func (w *wImpl) persistKeys() (int64, error) {
	sort.Sort(w.keys)
	kb := bb.New(size_page_key)      // page for keys
	bLeaves := make([]entries, 0, 0) // entries for last keys of pages, only offset relevant

	for _, e := range w.keys {
		if kb.Remaining() < 1+len(e.key) { // 1 byte for length, followed by key
			// TODO
		}
	}
}

var ErrKeyOverflow = errors.New("key length overflow")
var ErrDatOverflow = errors.New("data length overflow")

func (w *wImpl) Put(key string, dat []byte) (int, error) {
	if len(key) > MaxKeyLength {
		return 0, ErrKeyOverflow
	}
	if len(dat) > MaxDataLength {
		return 0, ErrDatOverflow
	}
	if len(dat) == 0 {
		return 0, nil
	}

	offset := w.cur
	n, err := w.put(dat)
	if err != nil || n != len(dat) {
		return n, err
	}

	w.keys = append(w.keys, entry{key, offset, int32(n)})
	return n, nil
}

func (w *wImpl) putBuf(buf bb.ByteBuffer) (int, error) {
	// make sure enough space available in buffer
	n := buf.Remaining()
	if n == 0 {
		return 0, nil
	}
	if w.buf.Remaining() < n {
		if err := fwc(w.buf, w.f); err != nil { // flip, write, clear
			return 0, err
		}
	}
	n, err := buf.WriteTo(w.buf)
	if err != nil {
		return n, err
	}
	w.cur += n
	return n, nil
}

func (w *wImpl) put(dat []byte) (int, error) {
	// make sure enough space available in buffer
	n := len(dat)
	if w.buf.Remaining() < n {
		if err := fwc(w.buf, w.f); err != nil { // flip, write, clear
			return 0, err
		}
	}
	w.buf.PutN(dat)
	w.cur += n
	return n, nil
}

func fwc(buf bb.ByteBuffer, w io.Writer) error {
	buf.Flip()
	for buf.HasRemaining() {
		if err := buf.WriteTo(w); err != nil {
			return err
		}
	}
	buf.Clear()
	return nil
}
