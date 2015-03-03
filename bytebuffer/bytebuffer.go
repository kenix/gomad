/*
	Package bytebuffer contains implementation of ByteBuffer similar to
	ByteBuffer in JDK.
*/
package bytebuffer

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	bi "encoding/binary"
	"errors"
	"fmt"
	"io"
	us "unsafe"
)

// ByteBuffer is a fix-sized buffer of bytes with Read and Write methods. Bytes
// are read and written as are, whereas wider unsigned integers are read and
// written according to this ByteBuffer's byte order, which defaults to the
// platform's endianness.
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

func (bb *ByteBuffer) String() string {
	return fmt.Sprintf("P%d;C%d;L%dM%d;%v;%v",
		bb.position, cap(bb.buf), bb.limit, bb.mark, bb.order, bb.buf)
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

const (
	BM00_08 = uint64(0x00000000000000FF) << (iota << 3)
	BM08_16
	BM16_24
	BM24_32
	BM32_40
	BM40_48
	BM48_56
	BM56_64
)

// ByteOrder is the underlying platform's endianness.
var ByteOrder bi.ByteOrder

func init() {
	// determine underlying platform's endianness
	ui := uint16(1)
	if 1 == (*[2]byte)(us.Pointer(&ui))[0] {
		ByteOrder = bi.LittleEndian
	} else {
		ByteOrder = bi.BigEndian
	}
}

// New creates and initializes a new ByteBuffer with the given fixed size.
func New(size int) *ByteBuffer {
	return &ByteBuffer{make([]byte, size, size), 0, size, -1, ByteOrder}
}

// Wrap creates a new ByteBuffer with content in bb
func Wrap(bb []byte) *ByteBuffer {
	return &ByteBuffer{bb, 0, len(bb), -1, ByteOrder}
}

// Order returns the underlying byte order
func (bb *ByteBuffer) Order() bi.ByteOrder { return bb.order }

// OrderTo sets the underlying byte order to o
func (bb *ByteBuffer) OrderTo(o bi.ByteOrder) *ByteBuffer {
	bb.order = o
	return bb
}

// Position returns the position from which Read and Write are performed.
func (bb *ByteBuffer) Position() int { return bb.position }

// PositionTo sets the position to p, panics if p is negative or greater than
// this ByteBuffer's limit. Mark is discarded if set and greater than p.
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
func (bb *ByteBuffer) Limit() int { return bb.limit }

// LimitTo sets the limit to l, panics if l is negative or greater than this
// ByteBuffer's capacity. Position is set to l if greater than l. Mark is
// discarded if set and greater than l.
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
func (bb *ByteBuffer) Capacity() int { return cap(bb.buf) }

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

func (bb *ByteBuffer) addr() us.Pointer { return us.Pointer(&bb.buf[0]) }

func (bb *ByteBuffer) pos() us.Pointer {
	return us.Pointer(uintptr(bb.addr()) + uintptr(bb.position))
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
func (bb *ByteBuffer) Remaining() int { return bb.limit - bb.position }

// HasRamaining denotes if there are any bytes between position and limit.
func (bb *ByteBuffer) HasRemaining() bool { return bb.Remaining() > 0 }

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
// and advances position by the len(s), panics if less than len(s) left for writing.
func (bb *ByteBuffer) PutN(s []byte) *ByteBuffer {
	length := len(s)
	if bb.position+length > cap(bb.buf) {
		panic(ErrOverflow)
	}
	C.memmove(bb.pos(), us.Pointer(&s[0]), C.size_t(length))
	bb.position += length
	return bb
}

// Get returns one byte at the current position and advances position by one,
// panics if no bytes left for reading.
func (bb *ByteBuffer) Get() byte {
	if bb.Remaining() < 1 {
		panic(ErrUnderflow)
	}
	var b byte
	C.memmove(us.Pointer(&b), bb.pos(), 1)
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
	if bb.order != ByteOrder {
		ui := i<<8 | i>>8
		C.memmove(bb.pos(), us.Pointer(&ui), 2)
	} else {
		C.memmove(bb.pos(), us.Pointer(&i), 2)
	}
	bb.position += 2
	return bb
}

// GetUint16 returns an uint16 from the current position and advances position
// by 2, panics if less than 2 bytes left for reading.
func (bb *ByteBuffer) GetUint16() uint16 {
	if bb.Remaining() < 2 {
		panic(ErrUnderflow)
	}
	var i uint16
	C.memmove(us.Pointer(&i), bb.pos(), 2)
	bb.position += 2
	if bb.order != ByteOrder {
		return i<<8 | i>>8
	}
	return i
}

// PutUint32 writes i into this ByteBuffer from the current position and
// advances position by 4, panics if less than 4 bytes left for writing.
func (bb *ByteBuffer) PutUint32(i uint32) *ByteBuffer {
	if bb.position+4 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	if bb.order != ByteOrder {
		ui := i<<24 | i<<8&uint32(BM16_24) | i>>8&uint32(BM08_16) | i>>24
		C.memmove(bb.pos(), us.Pointer(&ui), 4)
	} else {
		C.memmove(bb.pos(), us.Pointer(&i), 4)
	}
	bb.position += 4
	return bb
}

// GetUint32 returns an uint32 from the current position and advances position
// by 4, panics if less than 4 bytes left for reading.
func (bb *ByteBuffer) GetUint32() uint32 {
	if bb.Remaining() < 4 {
		panic(ErrUnderflow)
	}
	var i uint32
	C.memmove(us.Pointer(&i), bb.pos(), 4)
	bb.position += 4
	if bb.order != ByteOrder {
		return i<<24 | i<<8&uint32(BM16_24) | i>>8&uint32(BM08_16) | i>>24
	}
	return i
}

// PutUint64 writes i into this ByteBuffer from the current position and
// advances position by 8, panics if less than 8 bytes left for writing.
func (bb *ByteBuffer) PutUint64(i uint64) *ByteBuffer {
	if bb.position+8 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	if bb.order != ByteOrder {
		ui := i<<56 | i<<40&BM48_56 |
			i<<24&BM40_48 | i<<8&BM32_40 | i>>8&BM24_32 |
			i>>24&BM16_24 | i>>40&BM08_16 | i>>56
		C.memmove(bb.pos(), us.Pointer(&ui), 8)
	} else {
		C.memmove(bb.pos(), us.Pointer(&i), 8)
	}
	bb.position += 8
	return bb
}

// GetUint64 returns an uint64 from the current position and advances position
// by 8, panics if less than 8 bytes left for reading.
func (bb *ByteBuffer) GetUint64() uint64 {
	if bb.Remaining() < 8 {
		panic(ErrUnderflow)
	}
	var i uint64
	C.memmove(us.Pointer(&i), bb.pos(), 8)
	bb.position += 8
	if bb.order != ByteOrder {
		return i<<56 | i<<40&BM48_56 | i<<24&BM40_48 | i<<8&BM32_40 |
			i>>8&BM24_32 | i>>24&BM16_24 | i>>40&BM08_16 | i>>56
	}
	return i
}

// PutAll writes the given dat into this ByteBuffer from the current position
// advances position by the total length of dat, panics if dat type is not
// supported or this ByteBuffer cannot hold all of dat.
func (bb *ByteBuffer) PutAll(dat ...interface{}) *ByteBuffer {
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
