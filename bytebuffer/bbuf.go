package bytebuffer

import (
	"encoding/binary"
	"io"
)

type bbuf struct {
	buf      []byte
	position int
	limit    int
	mark     int
	order    binary.ByteOrder
}

func New(size int) ByteBuffer {
	return &bbuf{make([]byte, size, size), 0, size, -1, ByteOrder}
}

func Wrap(bb []byte) ByteBuffer {
	return &bbuf{bb, 0, len(bb), -1, ByteOrder}
}

func (bb *bbuf) Order() binary.ByteOrder {
	return bb.order
}

func (bb *bbuf) OrderTo(o binary.ByteOrder) ByteBuffer {
	bb.order = o
	return bb
}

func (bb *bbuf) Position() int {
	return bb.position
}

func (bb *bbuf) PositionTo(p int) ByteBuffer {
	if p < 0 || p > bb.limit {
		panic(ErrPosition)
	}
	bb.position = p
	if bb.mark > p {
		bb.mark = -1
	}
	return bb
}

func (bb *bbuf) Mark() ByteBuffer {
	bb.mark = bb.position
	return bb
}

func (bb *bbuf) Reset() ByteBuffer {
	if bb.mark < 0 {
		panic(ErrMark)
	}
	bb.position = bb.mark
	return bb
}

func (bb *bbuf) Limit() int {
	return bb.limit
}

func (bb *bbuf) LimitTo(l int) ByteBuffer {
	if l < 0 || l > cap(bb.buf) {
		panic(ErrLimit)
	}
	bb.limit = l
	if bb.position > l {
		bb.position = l
	}
	if bb.mark > l {
		bb.mark = -1
	}
	return bb
}

func (bb *bbuf) Capacity() int {
	return cap(bb.buf)
}

func (bb *bbuf) Clear() ByteBuffer {
	bb.position = 0
	bb.limit = cap(bb.buf)
	bb.mark = -1
	return bb
}

func (bb *bbuf) Compact() ByteBuffer {
	n := copy(bb.buf[0:], bb.buf[bb.position:bb.limit])
	bb.position = n
	bb.limit = cap(bb.buf)
	bb.mark = -1
	return bb
}

func (bb *bbuf) Flip() ByteBuffer {
	bb.limit = bb.position
	bb.position = 0
	bb.mark = -1
	return bb
}

func (bb *bbuf) Rewind() ByteBuffer {
	bb.position = 0
	bb.mark = -1
	return bb
}

func (bb *bbuf) Remaining() int {
	return bb.limit - bb.position
}

func (bb *bbuf) HasRemaining() bool {
	return bb.Remaining() > 0
}

func (bb *bbuf) Put(b byte) ByteBuffer {
	if bb.position >= cap(bb.buf) {
		panic(ErrOverflow)
	}
	bb.buf[bb.position] = b
	bb.position++
	return bb
}

func (bb *bbuf) PutN(s []byte) ByteBuffer {
	if bb.position+len(s) > cap(bb.buf) {
		panic(ErrOverflow)
	}
	copy(bb.buf[bb.position:], s)
	bb.position += len(s)
	return bb
}

func (bb *bbuf) Get() byte {
	if bb.Remaining() < 1 {
		panic(ErrUnderflow)
	}
	b := bb.buf[bb.position]
	bb.position++
	return b
}

func (bb *bbuf) GetN(n int) []byte {
	if bb.Remaining() < n {
		panic(ErrUnderflow)
	}
	s := make([]byte, n)
	copy(s, bb.buf[bb.position:])
	bb.position += n
	return s
}

func (bb *bbuf) PutUint16(i uint16) ByteBuffer {
	if bb.position+2 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	bb.order.PutUint16(bb.buf[bb.position:bb.position+2], i)
	bb.position += 2
	return bb
}

func (bb *bbuf) GetUint16() uint16 {
	if bb.position+2 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint16(bb.buf[bb.position : bb.position+2])
	bb.position += 2
	return ui
}

func (bb *bbuf) PutUint32(i uint32) ByteBuffer {
	if bb.position+4 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	bb.order.PutUint32(bb.buf[bb.position:bb.position+4], i)
	bb.position += 4
	return bb
}

func (bb *bbuf) GetUint32() uint32 {
	if bb.position+4 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint32(bb.buf[bb.position : bb.position+4])
	bb.position += 4
	return ui
}

func (bb *bbuf) PutUint64(i uint64) ByteBuffer {
	if bb.position+8 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	bb.order.PutUint64(bb.buf[bb.position:bb.position+8], i)
	bb.position += 8
	return bb
}

func (bb *bbuf) GetUint64() uint64 {
	if bb.position+8 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint64(bb.buf[bb.position : bb.position+8])
	bb.position += 8
	return ui
}

func (bb *bbuf) PutAll(dat ...interface{}) ByteBuffer {
	if bb.position+lenOf(dat...) > cap(bb.buf) {
		panic(ErrOverflow)
	}
	for _, v := range dat {
		switch t := v.(type) {
		case byte:
			bb.Put(t)
		case []byte:
			bb.PutN(t)
		case uint16:
			bb.PutUint16(t)
		case uint32:
			bb.PutUint32(t)
		case uint64:
			bb.PutUint64(t)
		case Binary:
			bb.PutN(t.Bytes())
		}
	}
	return bb
}

func lenOf(dat ...interface{}) int {
	n := 0
	for _, v := range dat {
		switch t := v.(type) {
		case byte:
			n += 1
		case []byte:
			n += len(t)
		case uint16:
			n += 2
		case uint32:
			n += 4
		case uint64:
			n += 8
		case Binary:
			n += len(t.Bytes())
		default:
			panic(ErrType)
		}
	}
	return n
}

func (bb *bbuf) Read(s []byte) (int, error) {
	n := copy(s, bb.buf[bb.position:bb.limit])
	bb.position += n
	return n, nil
}

func (bb *bbuf) Write(s []byte) (int, error) {
	n := copy(bb.buf[bb.position:], s)
	bb.position += n
	return n, nil
}

func (bb *bbuf) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(bb.buf[bb.position:])
	bb.position += n
	return int64(n), err
}

func (bb *bbuf) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(bb.buf[bb.position:bb.limit])
	bb.position += n
	return int64(n), err
}
