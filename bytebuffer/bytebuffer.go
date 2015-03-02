/*
	Package bytebuffer contains implementation of ByteBuffer similar to Java.
*/
package bytebuffer

/*
#include <stdlib.h>
#include <string.h>

char at(char* addr){
	return *addr;
}
*/
import "C"
import (
	bi "encoding/binary"
	"errors"
	"io"
	us "unsafe"
)

// ByteBuffer is a fix-sized buffer of bytes with Read and Write methods. Bytes
// are read and written as are, whereas wider unsigned integers are read and
// written according to this ByteBuffer's byte order.
//
// Reading operations are performed from position to limit.
// Writing operations are performed from position to capacity.
type ByteBuffer struct {
	buf      []byte
	position int
	limit    int
	mark     int
	order    bi.ByteOrder
}

// Binary denotes something that has bytes
type Binary interface {
	Bytes() []byte
}

var ErrOverflow = errors.New("bytebuffer: overflow")
var ErrUnderflow = errors.New("bytebuffer: underflow")
var ErrPosition = errors.New("bytebuffer: invalid position")
var ErrLimit = errors.New("bytebuffer: invalid limit")
var ErrMark = errors.New("bytebuffer: mark undefined")
var ErrType = errors.New("bytebuffer: unsupported type")

// New creates and initializes a new ByteBuffer with the given fixed size.
func New(size int) *ByteBuffer {
	return &ByteBuffer{make([]byte, size, size), 0, size, -1, bi.BigEndian}
}

// Wrap creates a new ByteBuffer with content in bb
func Wrap(bb []byte) *ByteBuffer {
	return &ByteBuffer{bb, 0, len(bb), -1, bi.BigEndian}
}

// Order returns the underlying byte order
func (bb *ByteBuffer) Order() bi.ByteOrder {
	return bb.order
}

// OrderTo sets the underlying byte order to o
func (bb *ByteBuffer) OrderTo(o bi.ByteOrder) *ByteBuffer {
	bb.order = o
	return bb
}

// Position returns the position from which Read and Write are performed.
func (bb *ByteBuffer) Position() int {
	return bb.position
}

// PositionTo sets the position to p, panics if p is negative or greater than
// this ByteBuffer's limit.
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

// Mark marks this ByteBuffer at its current position.
func (bb *ByteBuffer) Mark() *ByteBuffer {
	bb.mark = bb.position
	return bb
}

// Reset sets this ByteBuffer's current position to its mark, panics if mark
// was not set.
func (bb *ByteBuffer) Reset() *ByteBuffer {
	if bb.mark < 0 {
		panic(ErrMark)
	}
	bb.position = bb.mark
	return bb
}

// Limit returns this ByteBuffer's limit.
func (bb *ByteBuffer) Limit() int {
	return bb.limit
}

// LimitTo sets the limit to l, panics if l is negative or greater than this
// ByteBuffer's capacity.
func (bb *ByteBuffer) LimitTo(l int) *ByteBuffer {
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

// Capacity returns the capacity of this ByteBuffer. It determines how many
// bytes this ByteBuffer can hold.
func (bb *ByteBuffer) Capacity() int {
	return cap(bb.buf)
}

// Clear sets this ByteBuffer's cursor states back as it was created. Its
// content is not touched.
func (bb *ByteBuffer) Clear() *ByteBuffer {
	bb.position = 0
	bb.limit = cap(bb.buf)
	bb.mark = -1
	return bb
}

// Compact moves the n bytes between position and limit to the beginning of this
// ByteBuffer. Sets position to n, limit to capacity and discards mark.
func (bb *ByteBuffer) Compact() *ByteBuffer {
	n := bb.limit - bb.position
	if n > 0 {
		C.memmove(bb.addr(), bb.pos(), C.size_t(n))
	}
	bb.position = n
	bb.limit = cap(bb.buf)
	bb.mark = -1
	return bb
}

func (bb *ByteBuffer) addr() us.Pointer {
	return us.Pointer(&bb.buf[0])
}

func (bb *ByteBuffer) pos() us.Pointer {
	return us.Pointer(uintptr(bb.addr()) + uintptr(bb.position))
}

func (bb *ByteBuffer) asCharA() *C.char {
	return (*C.char)(bb.pos())
}

// Flip sets limit to position, position to 0 and discards mark, readies this
// ByteBuffer for reading.
func (bb *ByteBuffer) Flip() *ByteBuffer {
	bb.limit = bb.position
	bb.position = 0
	bb.mark = -1
	return bb
}

// Rewind sets position to 0 and discards mark, readies this ByteBuffer for
// re-read.
func (bb *ByteBuffer) Rewind() *ByteBuffer {
	bb.position = 0
	bb.mark = -1
	return bb
}

// Remaining returns the number of bytes between position and limit.
func (bb *ByteBuffer) Remaining() int {
	return bb.limit - bb.position
}

// HasRamaining denotes if there are any bytes between position and limit.
func (bb *ByteBuffer) HasRemaining() bool {
	return bb.Remaining() > 0
}

// Put writes one byte into this ByteBuffer at the current position and advances
// position by one, panics if capacity is reached.
func (bb *ByteBuffer) Put(b byte) *ByteBuffer {
	if bb.position >= cap(bb.buf) {
		panic(ErrOverflow)
	}
	C.memset(bb.pos(), C.int(b), 1)
	bb.position++
	return bb
}

// PutN writes len(s) bytes in s into this ByteBuffer from the current position
// and advances position by the len(s), panics if less tha len(s) left for writing.
func (bb *ByteBuffer) PutN(s []byte) *ByteBuffer {
	l := len(s)
	if bb.position+l > cap(bb.buf) {
		panic(ErrOverflow)
	}
	C.memmove(bb.pos(), us.Pointer(&s[0]), C.size_t(l))
	bb.position += l
	return bb
}

// Get returns one byte at the current position and advances position by one,
// panics if no bytes left for reading.
func (bb *ByteBuffer) Get() byte {
	if bb.Remaining() < 1 {
		panic(ErrUnderflow)
	}
	b := byte(C.at(bb.asCharA()))
	bb.position += 1
	return b
}

// GetN returns a slice of n bytes from the current position and advances
// position by n, panics if less than n bytes left for reading.
func (bb *ByteBuffer) GetN(n int) []byte {
	if bb.Remaining() < n {
		panic(ErrUnderflow)
	}
	s := make([]byte, n)
	C.memmove(us.Pointer(&s[0]), bb.pos(), C.size_t(n))
	bb.position += n
	return s
}

// PutUint16 writes i into this ByteBuffer from the current position and
// advances position by 2, panics if less than 2 bytes left for writing.
func (bb *ByteBuffer) PutUint16(i uint16) *ByteBuffer {
	if bb.position+2 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	if bb.order == bi.LittleEndian {
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
	} else {
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
	}
	return bb
}

// GetUint16 returns an uint16 from the current position and advances position
// by 2, panics if less than 2 bytes left for reading.
func (bb *ByteBuffer) GetUint16() uint16 {
	if bb.position+2 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := uint16(0)
	if bb.order == bi.LittleEndian {
		ui |= uint16(C.at(bb.asCharA()))
		bb.position += 1
		ui |= uint16(C.at(bb.asCharA())) << 8
		bb.position += 1
	} else {
		ui |= uint16(C.at(bb.asCharA())) << 8
		bb.position += 1
		ui |= uint16(C.at(bb.asCharA()))
		bb.position += 1
	}
	return ui
}

// PutUint32 writes i into this ByteBuffer from the current position and
// advances position by 4, panics if less than 4 bytes left for writing.
func (bb *ByteBuffer) PutUint32(i uint32) *ByteBuffer {
	if bb.position+4 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	if bb.order == bi.LittleEndian {
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>16)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>24)), 1)
		bb.position += 1
	} else {
		C.memset(bb.pos(), C.int(byte(i>>24)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>16)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
	}
	return bb
}

// GetUint32 returns an uint32 from the current position and advances position
// by 4, panics if less than 4 bytes left for reading.
func (bb *ByteBuffer) GetUint32() uint32 {
	if bb.position+4 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := uint32(0)
	if bb.order == bi.LittleEndian {
		ui |= uint32(C.at(bb.asCharA()))
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA())) << 8
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA())) << 16
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA())) << 24
		bb.position += 1
	} else {
		ui |= uint32(C.at(bb.asCharA())) << 24
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA())) << 16
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA())) << 8
		bb.position += 1
		ui |= uint32(C.at(bb.asCharA()))
		bb.position += 1
	}
	return ui
}

// PutUint64 writes i into this ByteBuffer from the current position and
// advances position by 8, panics if less than 8 bytes left for writing.
func (bb *ByteBuffer) PutUint64(i uint64) *ByteBuffer {
	if bb.position+8 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	if bb.order == bi.LittleEndian {
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>16)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>24)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>32)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>40)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>48)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>56)), 1)
		bb.position += 1
	} else {
		C.memset(bb.pos(), C.int(byte(i>>56)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>48)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>40)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>32)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>24)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>16)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i>>8)), 1)
		bb.position += 1
		C.memset(bb.pos(), C.int(byte(i)), 1)
		bb.position += 1
	}
	return bb
}

// GetUint64 returns an uint64 from the current position and advances position
// by 8, panics if less than 8 bytes left for reading.
func (bb *ByteBuffer) GetUint64() uint64 {
	if bb.position+8 > bb.limit {
		panic(ErrUnderflow)
	}
	ui := uint64(0)
	if bb.order == bi.LittleEndian {
		ui |= uint64(C.at(bb.asCharA()))
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 8
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 16
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 24
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 32
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 40
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 48
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 56
		bb.position += 1
	} else {
		ui |= uint64(C.at(bb.asCharA())) << 56
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 48
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 40
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 32
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 24
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 16
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA())) << 8
		bb.position += 1
		ui |= uint64(C.at(bb.asCharA()))
		bb.position += 1
	}
	return ui
}

// PutAll writes the given dat into this ByteBuffer from the current position
// advances position by the total length of dat, panics if dat type is not
// supported or this ByteBuffer cannot hold all of dat, which can result in
// partial write as for now.
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

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (bb *ByteBuffer) Read(s []byte) int {
	n := min(bb.Remaining(), len(s))
	C.memmove(us.Pointer(&s[0]), bb.pos(), C.size_t(n))
	bb.position += n
	return n
}

func (bb *ByteBuffer) Write(s []byte) int {
	n := min(cap(bb.buf)-bb.position, len(s))
	C.memmove(bb.pos(), us.Pointer(&s[0]), C.size_t(n))
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
