package sdb

import (
	"errors"
	"io"
	"os"
	"sort"

	bb "github.com/kenix/gomad/bytebuffer"
)

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
	return &wImpl{f, bb.New(MaxDataLength), make([]*entry, 0, 0), 0}, nil
}

func (w *wImpl) Underlying() string {
	return w.f.Name()
}

func (w *wImpl) Close() error {
	defer w.f.Close() // make sure underlying file is always closed

	bls, err := w.persistKeys() // B*-tree level 1 leaves start position
	if err != nil {
		return err
	}
	w.buf.PutUint64(uint64(bls))
	if err := fwc(w.buf, w.f); err != nil {
		return err
	}
	return nil
}

func (w *wImpl) persistKeys() (int64, error) {
	sort.Sort(w.keys)
	kb := bb.New(size_page_key)     // page for keys
	bLeaves := make([]*entry, 0, 0) // entries for last keys of pages, only offset relevant

	offset := w.cur
	for _, e := range w.keys {
		if kb.Remaining() < 1+len(e.key)+8 { // 1 byte for key length, key, offset and data length
			if _, err := w.fwcBuf(kb); err != nil {
				return 0, err
			}
			bLeaves = append(bLeaves, &entry{e.key, offset, 0})
			offset = w.cur
		}
		kb.Put(byte(len(e.key))).PutN([]byte(e.key))
		kb.PutUint64(uint64(e.offset<<dataLengthBits) | uint64(e.length))
	}

	if kb.Position() > 0 {
		if _, err := w.fwcBuf(kb); err != nil {
			return 0, err
		}
	}

	bls := w.cur
	if bls > offset { // data exists
		bLeaves = append(bLeaves, &entry{"", offset, 0})
	}

	if len(bLeaves) > 0 {
		pkb := bb.New(1 + MaxKeyLength + 8)
		for _, b := range bLeaves {
			pkb.PutUint64(uint64(b.offset))
			if len(b.key) > 0 {
				pkb.Put(byte(len(b.key))).PutN([]byte(b.key))
			}
			if _, err := w.fwcBuf(pkb); err != nil {
				return bls, err
			}
		}
	}

	return bls, nil
}

var ErrKeyOverflow = errors.New("key length overflow")
var ErrDatOverflow = errors.New("data length overflow")
var ErrPartialWrite = errors.New("partial write")

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

	w.keys = append(w.keys, &entry{key, offset, int32(n)})
	return n, nil
}

func (w *wImpl) fwcBuf(buf bb.ByteBuffer) (int, error) {
	buf.Flip()
	defer buf.Clear()
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
	nw, err := buf.WriteTo(w.buf)
	if err != nil {
		return int(nw), err
	}
	if int(nw) != n {
		return int(nw), ErrPartialWrite
	}
	w.cur += int64(n)
	return n, nil
}

func (w *wImpl) put(dat []byte) (int, error) {
	// make sure enough space available in buffer
	n := len(dat)
	if n == 0 {
		return 0, nil
	}
	if w.buf.Remaining() < n {
		if err := fwc(w.buf, w.f); err != nil { // flip, write, clear
			return 0, err
		}
	}
	w.buf.PutN(dat)
	w.cur += int64(n)
	return n, nil
}

func fwc(buf bb.ByteBuffer, w io.Writer) error {
	buf.Flip()
	defer buf.Clear()
	for buf.HasRemaining() {
		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}
	return nil
}
