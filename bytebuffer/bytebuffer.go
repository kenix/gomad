package bytebuffer

import (
	"encoding/binary"
	"errors"
	"io"
)

type ByteBuffer struct {
	buf      []byte
	capacity int
	position int
	limit    int
	mark     int
	order    binary.ByteOrder
}

type Binary interface {
	Bytes() []byte
}

var ErrOverflow = errors.New("bytebuffer: overflow")
var ErrUnderflow = errors.New("bytebuffer: underflow")
var ErrPosition = errors.New("bytebuffer: invalid position")
var ErrLimit = errors.New("bytebuffer: invalid limit")
var ErrMark = errors.New("bytebuffer: mark undefined")
var ErrType = errors.New("bytebuffer: unsupported type")

func New(size int) *ByteBuffer {
	return &ByteBuffer{make([]byte, size, size),
		size, 0, size, -1, binary.BigEndian}
}

func Wrap(bb []byte) *ByteBuffer {
	return &ByteBuffer{bb, len(bb), 0, len(bb), -1, binary.BigEndian}
}

func (bb *ByteBuffer) Order() binary.ByteOrder {
	return bb.order
}

func (bb *ByteBuffer) OrderTo(o binary.ByteOrder) *ByteBuffer {
	bb.order = o
	return bb
}

func (bb *ByteBuffer) Position() int {
	return bb.position
}

func (bb *ByteBuffer) PositionTo(p int) *ByteBuffer {
	if p < 0 || p > bb.limit {
		panic(ErrPosition)
	}
	bb.position = p
	if bb.mark > p {
		bb.mark = -1
	}
	return bb
}

func (bb *ByteBuffer) Mark() *ByteBuffer {
	bb.mark = bb.position
	return bb
}

func (bb *ByteBuffer) Reset() *ByteBuffer {
	if bb.mark < 0 {
		panic(ErrMark)
	}
	bb.position = bb.mark
	return bb
}

func (bb *ByteBuffer) Limit() int {
	return bb.limit
}

func (bb *ByteBuffer) LimitTo(l int) *ByteBuffer {
	if l < 0 || l > bb.capacity {
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

func (bb *ByteBuffer) Capacity() int {
	return bb.capacity
}

func (bb *ByteBuffer) Clear() *ByteBuffer {
	bb.position = 0
	bb.limit = bb.capacity
	bb.mark = -1
	return bb
}

func (bb *ByteBuffer) Compact() *ByteBuffer {
	n := copy(bb.buf[0:], bb.buf[bb.position:bb.limit])
	bb.position = n
	bb.limit = bb.capacity
	bb.mark = -1
	return bb
}

func (bb *ByteBuffer) Flip() *ByteBuffer {
	bb.limit = bb.position
	bb.position = 0
	bb.mark = -1
	return bb
}

func (bb *ByteBuffer) Rewind() *ByteBuffer {
	bb.position = 0
	bb.mark = -1
	return bb
}

func (bb *ByteBuffer) Remaining() int {
	return bb.limit - bb.position
}

func (bb *ByteBuffer) HasRemaining() bool {
	return bb.Remaining() > 0
}

func (bb *ByteBuffer) Put(b byte) *ByteBuffer {
	if bb.position >= bb.capacity {
		panic(ErrOverflow)
	}
	bb.buf[bb.position] = b
	bb.position++
	return bb
}

func (bb *ByteBuffer) PutN(s []byte) *ByteBuffer {
	if bb.position+len(s) > bb.capacity {
		panic(ErrOverflow)
	}
	copy(bb.buf[bb.position:], s)
	bb.position += len(s)
	return bb
}

func (bb *ByteBuffer) Get() byte {
	if bb.Remaining() < 1 {
		panic(ErrUnderflow)
	}
	b := bb.buf[bb.position]
	bb.position++
	return b
}

func (bb *ByteBuffer) GetN(n int) []byte {
	if bb.Remaining() < n {
		panic(ErrUnderflow)
	}
	s := make([]byte, n)
	copy(s, bb.buf[bb.position:])
	bb.position += n
	return s
}

func (bb *ByteBuffer) PutUint16(i uint16) *ByteBuffer {
	if bb.position+2 > bb.capacity {
		panic(ErrOverflow)
	}
	bb.order.PutUint16(bb.buf[bb.position:bb.position+2], i)
	bb.position += 2
	return bb
}

func (bb *ByteBuffer) GetUint16() uint16 {
	if bb.position+2 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint16(bb.buf[bb.position : bb.position+2])
	bb.position += 2
	return ui
}

func (bb *ByteBuffer) PutUint32(i uint32) *ByteBuffer {
	if bb.position+4 > bb.capacity {
		panic(ErrOverflow)
	}
	bb.order.PutUint32(bb.buf[bb.position:bb.position+4], i)
	bb.position += 4
	return bb
}

func (bb *ByteBuffer) GetUint32() uint32 {
	if bb.position+4 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint32(bb.buf[bb.position : bb.position+4])
	bb.position += 4
	return ui
}

func (bb *ByteBuffer) PutUint64(i uint64) *ByteBuffer {
	if bb.position+8 > bb.capacity {
		panic(ErrOverflow)
	}
	bb.order.PutUint64(bb.buf[bb.position:bb.position+8], i)
	bb.position += 8
	return bb
}

func (bb *ByteBuffer) GetUint64() uint64 {
	if bb.position+8 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := bb.order.Uint64(bb.buf[bb.position : bb.position+8])
	bb.position += 8
	return ui
}

func (bb *ByteBuffer) PutAll(dat ...interface{}) *ByteBuffer {
	for _, val := range dat {
		switch t := val.(type) {
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
		default:
			panic(ErrType)
		}
	}
	return bb
}

func (bb *ByteBuffer) Read(s []byte) int {
	n := copy(s, bb.buf[bb.position:bb.limit])
	bb.position += n
	return n
}

func (bb *ByteBuffer) Write(s []byte) int {
	n := copy(bb.buf[bb.position:], s)
	bb.position += n
	return n
}

func (bb *ByteBuffer) ReadFrom(r io.Reader) (int, error) {
	n, err := r.Read(bb.buf[bb.position:])
	bb.position += n
	return n, err
}

func (bb *ByteBuffer) WriteTo(w io.Writer) (int, error) {
	n, err := w.Write(bb.buf[bb.position:bb.limit])
	bb.position += n
	return n, err
}
