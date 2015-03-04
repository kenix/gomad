/*
	Package bytebuffer contains implementation of ByteBuffer similar to
	ByteBuffer in JDK.
*/
package bytebuffer

import (
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

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
var ByteOrder binary.ByteOrder

func init() {
	// determine underlying platform's endianness
	ui := uint16(1)
	if 1 == (*[2]byte)(unsafe.Pointer(&ui))[0] {
		ByteOrder = binary.LittleEndian
	} else {
		ByteOrder = binary.BigEndian
	}
}

// ByteBuffer is a fix-sized buffer of bytes with Read and Write methods. Bytes
// are read and written as are, whereas wider unsigned integers are read and
// written according to this ByteBuffer's byte order, which defaults to the
// platform's endianness.
//
// Reading operations are performed from position to limit.
// Writing operations are performed from position to capacity.
type ByteBuffer interface {
	// Order returns the underlying byte order
	Order() binary.ByteOrder

	// OrderTo sets the underlying byte order to o
	OrderTo(o binary.ByteOrder) ByteBuffer

	// Position returns the position from which Read and Write are performed.
	Position() int

	// PositionTo sets the position to p, panics if p is negative or greater than
	// this ByteBuffer's limit. Mark is discarded if set and greater than p.
	PositionTo(p int) ByteBuffer

	// Mark marks this ByteBuffer at its current position.
	Mark() ByteBuffer

	// Reset sets this ByteBuffer's current position to its mark, panics if mark
	// was not set.
	Reset() ByteBuffer

	// Limit returns this ByteBuffer's limit.
	Limit() int

	// LimitTo sets the limit to l, panics if l is negative or greater than this
	// ByteBuffer's capacity. Position is set to l if greater than l. Mark is
	// discarded if set and greater than l.
	LimitTo(l int) ByteBuffer

	// Capacity returns the capacity of this ByteBuffer. It determines how many
	// bytes this ByteBuffer can hold.
	Capacity() int

	// Clear sets this ByteBuffer's cursor states back as it was created. Its
	// content is not touched.
	Clear() ByteBuffer

	// Compact moves the n bytes between position and limit to the beginning of this
	// ByteBuffer. Sets position to n, limit to capacity and discards mark.
	Compact() ByteBuffer

	// Flip sets limit to position, position to 0 and discards mark, readies this
	// ByteBuffer for reading.
	Flip() ByteBuffer

	// Rewind sets position to 0 and discards mark, readies this ByteBuffer for
	// re-read.
	Rewind() ByteBuffer

	// Remaining returns the number of bytes between position and limit.
	Remaining() int

	// HasRamaining denotes if there are any bytes between position and limit.
	HasRemaining() bool

	// Put writes one byte into this ByteBuffer at the current position and advances
	// position by one, panics if capacity is reached.
	Put(b byte) ByteBuffer

	// PutN writes len(s) bytes in s into this ByteBuffer from the current position
	// and advances position by the len(s), panics if less than len(s) left for writing.
	PutN(s []byte) ByteBuffer

	// Get returns one byte at the current position and advances position by one,
	// panics if no bytes left for reading.
	Get() byte

	// GetN returns a slice of n bytes from the current position and advances
	// position by n, panics if less than n bytes left for reading.
	GetN(n int) []byte

	// PutUint16 writes i into this ByteBuffer from the current position and
	// advances position by 2, panics if less than 2 bytes left for writing.
	PutUint16(i uint16) ByteBuffer

	// GetUint16 returns an uint16 from the current position and advances position
	// by 2, panics if less than 2 bytes left for reading.
	GetUint16() uint16

	// PutUint32 writes i into this ByteBuffer from the current position and
	// advances position by 4, panics if less than 4 bytes left for writing.
	PutUint32(i uint32) ByteBuffer

	// GetUint32 returns an uint32 from the current position and advances position
	// by 4, panics if less than 4 bytes left for reading.
	GetUint32() uint32

	// PutUint64 writes i into this ByteBuffer from the current position and
	// advances position by 8, panics if less than 8 bytes left for writing.
	PutUint64(i uint64) ByteBuffer

	// GetUint64 returns an uint64 from the current position and advances position
	// by 8, panics if less than 8 bytes left for reading.
	GetUint64() uint64

	// PutAll writes the given dat into this ByteBuffer from the current position
	// advances position by the total length of dat, panics if dat type is not
	// supported or this ByteBuffer cannot hold all of dat.
	PutAll(dat ...interface{}) ByteBuffer

	io.Reader
	io.Writer
	io.ReaderFrom
	io.WriterTo
}
