package sdb

import (
	"io"
	"math"
)

const (
	dataLengthBits = 24
	dataLengthMask = 1<<24 - 1
	// Currently supported maximum key length 255
	MaxKeyLength = math.MaxUint8
	// Currently supported maximum data length 16M
	MaxDataLength = 1 << dataLengthBits
	size_page_key = 1 << 12 // 4K
)

// Underlyer denotes something that has a underlying file on disk.
type Underlyer interface {
	// Underlying returns the path of the underlying file on disk.
	Underlying() string
	io.Closer
}

/*
Writer stores the given dat under key in the underlying file on disk. The last write
of storing data with same key multiple times will win. Data written is buffered.
After closing this writer the data and indices are persisted in the underlying file
on disk.

Typical usage is to create the writer, write the key and data pairs in a loop, then
close the writer.
*/
type Writer interface {
	// Put writes the given dat under key in the underlying file on disk, returns
	// bytes written and error. err will be nil if the write is successful.
	Put(key string, dat []byte) (n int, err error)
	// Underlying returns the path of the underlying file on disk.
	Underlyer
}

/*
Reader provides the interface for querying data stored in the underlying file on disk
for a given key.

Typical usage is to create the reader and keep it open as long as the underlying
file doesn't change. Providing querying for keys. After service close the reader.
*/
type Reader interface {
	// Get reads the data for the given key and returns it in a byte slice. error
	// will be not nil if the read or query is not successful.
	Get(key string) ([]byte, error)
	Underlyer
}
