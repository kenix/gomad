package bytebuffer

import (
	"bytes"
	bi "encoding/binary"
	"fmt"
	"math"
	"testing"
	us "unsafe"
)

func BenchmarkPut_S(b *testing.B) {
	s := []byte(nil)
	for i := 0; i < b.N; i++ {
		s = append(s, byte(math.MaxUint8))
	}
	b.ReportAllocs()
}

func BenchmarkPut(b *testing.B) {
	benchmarkPut(b, New(1<<20), byte(math.MaxUint8))
}

func benchmarkPut(b *testing.B, bb ByteBuffer, d interface{}) {
	switch t := d.(type) {
	case byte:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Clear()
			}
			bb.Put(t)
		}
	case uint16:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Clear()
			}
			bb.PutUint16(t)
		}
	case uint32:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Clear()
			}
			bb.PutUint32(t)
			b.SetBytes(4)
		}
	case uint64:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Clear()
			}
			bb.PutUint64(t)
			b.SetBytes(8)
		}
	case []byte:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Clear()
			}
			bb.PutN(t)
			b.SetBytes(int64(len(t)))
		}
	default:
		b.Errorf("no support for %T\n", t)
	}
	b.ReportAllocs()
}

func BenchmarkGet(b *testing.B) {
	benchmarkGet(b, New(1<<20), byte(math.MaxUint8))
}

func benchmarkGet(b *testing.B, bb ByteBuffer, d interface{}) {
	n := 0
	for bb.HasRemaining() {
		switch t := d.(type) {
		case byte:
			bb.Put(t)
		case uint16:
			bb.PutUint16(t)
		case uint32:
			bb.PutUint32(t)
		case uint64:
			bb.PutUint64(t)
		case []byte:
			bb.PutN(t)
			n = len(t)
		default:
			b.Errorf("no support for %T\n", t)
		}
	}
	bb.Flip()

	b.ResetTimer()
	switch t := d.(type) {
	case byte:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Flip()
			}
			bb.Get()
		}
	case uint16:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Flip()
			}
			bb.GetUint16()
		}
	case uint32:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Flip()
			}
			bb.GetUint32()
			b.SetBytes(4)
		}
	case uint64:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Flip()
			}
			bb.GetUint64()
			b.SetBytes(8)
		}
	case []byte:
		for i := 0; i < b.N; i++ {
			if !bb.HasRemaining() {
				bb.Flip()
			}
			bb.GetN(n)
			b.SetBytes(int64(n))
		}
	default:
		b.Errorf("no support for %T\n", t)
	}
	b.ReportAllocs()
}

func benchmarkRead(b *testing.B, bb ByteBuffer, size int) {
	for bb.HasRemaining() {
		bb.Put(math.MaxUint8)
	}
	bb.Flip()

	s := make([]byte, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !bb.HasRemaining() {
			bb.Flip()
		}
		bb.Read(s)
		b.SetBytes(int64(size))
	}
	b.ReportAllocs()
}

func BenchmarkPutUint16(b *testing.B) {
	benchmarkPut(b, New(1<<20), uint16(math.MaxUint16))
}

func BenchmarkGetUint16(b *testing.B) {
	benchmarkGet(b, New(1<<20), uint16(math.MaxUint16))
}

func BenchmarkPutUint32(b *testing.B) {
	benchmarkPut(b, New(1<<20), uint32(math.MaxUint32))
}

func BenchmarkGetUint32(b *testing.B) {
	benchmarkGet(b, New(1<<20), uint32(math.MaxUint32))
}

func BenchmarkPutUint64(b *testing.B) {
	benchmarkPut(b, New(1<<20), uint64(math.MaxUint64))
}

func BenchmarkGetUint64(b *testing.B) {
	benchmarkGet(b, New(1<<20), uint64(math.MaxUint64))
}

func makeSlice(size int) []byte {
	s := make([]byte, size)
	for i := range s {
		s[i] = math.MaxUint8
	}
	return s
}

func BenchmarkPut16(b *testing.B) {
	benchmarkPut(b, New(1<<20), makeSlice(16))
}

func BenchmarkPut32(b *testing.B) {
	benchmarkPut(b, New(1<<20), makeSlice(32))
}

func BenchmarkPut64(b *testing.B) {
	benchmarkPut(b, New(1<<20), makeSlice(64))
}

func BenchmarkGet16(b *testing.B) {
	benchmarkGet(b, New(1<<20), makeSlice(16))
}

func BenchmarkGet32(b *testing.B) {
	benchmarkGet(b, New(1<<20), makeSlice(32))
}

func BenchmarkGet64(b *testing.B) {
	benchmarkGet(b, New(1<<20), makeSlice(64))
}

func BenchmarkRead16(b *testing.B) {
	benchmarkRead(b, New(1<<20), 16)
}

func BenchmarkRead32(b *testing.B) {
	benchmarkRead(b, New(1<<20), 32)
}

func BenchmarkRead64(b *testing.B) {
	benchmarkRead(b, New(1<<20), 64)
}

func TestPutAll(t *testing.T) {
	bb := New(32)
	bb.PutAll(byte(1), uint16(2), uint32(3), uint64(4), // 15 bytes
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},               // 9 bytes
		bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8})) // 8 bytes
	checkErrCase(t, ErrOverflow, func() { bb.PutAll([]byte{1, 2}) })
	checkCursors(t, bb, 32, 32, 32)
	bb.Flip()

	if bb.Get() != 1 || bb.GetUint16() != 2 || bb.GetUint32() != 3 || bb.GetUint64() != 4 {
		t.Error("numbers: wanted 1,2,3,4")
	}
	s := bb.GetN(9)
	if len(s) != 9 || fmt.Sprintf("%v", s) != "[1 2 3 4 5 6 7 8 9]" {
		t.Error("bytes: wanted [1 2 3 4 5 6 7 8 9]")
	}
	buf := bytes.NewBuffer(bb.GetN(8))
	if buf.Len() != 8 || fmt.Sprintf("%v", buf.Bytes()) != "[1 2 3 4 5 6 7 8]" {
		t.Error("Binary: wanted [1 2 3 4 5 6 7 8]")
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })

	bb.Clear()
	checkErrCase(t, ErrType, func() { bb.PutAll(int16(2)) })
	checkErrCase(t, ErrType, func() { bb.PutAll(float32(3.14)) })
}

func TestGetN(t *testing.T) {
	bb := New(8)
	bb.Put(1)
	dat := []byte{2, 3, 4, 5}
	bb.PutN(dat)
	bb.Flip()
	checkErrCase(t, ErrUnderflow, func() { bb.GetN(16) })
	dat2 := bb.GetN(4)
	for i, d := range dat2 {
		if i+1 != int(d) {
			t.Errorf("Byte: wanted %d, got %d", i+1, d)
		}
	}
	if bb.Get() != 5 {
		t.FailNow()
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })
}

func TestPutN(t *testing.T) {
	bb := New(8)
	bb.Put(1)
	dat := []byte{2, 3, 4, 5}
	bb.PutN(dat)
	checkCursors(t, bb, 8, 5, 8)
	checkErrCase(t, ErrOverflow, func() { bb.PutN(dat) })
	bb.Flip()
	bb.Get()
	for _, d := range dat {
		b := bb.Get()
		if b != d {
			t.Errorf("Byte: wanted %d, got %d", d, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })
}

func TestWithWriter(t *testing.T) {
	bb := New(16)
	bb.Put(1)
	bb.Write([]byte{2, 3, 4, 5})
	bb.Put(6)
	bb.Flip()

	var buf bytes.Buffer
	bb.Get()
	n, err := bb.WriteTo(&buf)
	if n != 5 || err != nil {
		t.FailNow()
	}
	checkCursors(t, bb, 16, 6, 6)

	s := buf.Bytes()
	for i := 2; i < 7; i++ {
		if s[i-2] != byte(i) {
			t.Errorf("Byte: wanted %d, got %d", i, s[i-2])
		}
	}
}

func TestWithReader(t *testing.T) {
	bb := New(16)
	bb.Put(1)
	n, err := bb.ReadFrom(bytes.NewReader([]byte{2, 3, 4, 5}))
	if n != 4 || err != nil {
		t.FailNow()
	}
	bb.Put(6)
	bb.Flip()
	for i := 1; i <= 6; i++ {
		b := bb.Get()
		if b != byte(i) {
			t.Errorf("Element: wanted %d, got %d", i, b)
		}
	}
}

func TestBytesWrite(t *testing.T) {
	bb := New(16)
	bb.Put(1)
	n, err := bb.Write([]byte{2, 3, 4, 5})
	if err != nil {
		t.FailNow()
	}
	if n != 4 {
		t.FailNow()
	}
	bb.Put(6)
	bb.Flip()
	for i := 1; i <= 6; i++ {
		b := bb.Get()
		if b != byte(i) {
			t.Errorf("Element: wanted %d, got %d", i, b)
		}
	}
}

func TestBytesRead(t *testing.T) {
	bb := New(16)
	bb.Put(1)
	bb.Write([]byte{2, 3, 4, 5})
	bb.Put(6)
	bb.Flip()
	s := make([]byte, 4)

	bb.Read(s[:1])
	if s[0] != 1 {
		t.Errorf("Byte: wanted 1, got %d", s[0])
	}
	bb.Read(s)
	for i := 2; i < 6; i++ {
		if s[i-2] != byte(i) {
			t.Errorf("Byte: wanted %d, got %d", i, s[i-2])
		}
	}
	bb.Read(s[3:])
	checkCursors(t, bb, 16, 6, 6)
	if s[3] != 6 {
		t.Errorf("Byte: wanted %d, got %d", 6, s[3])
	}
}

func TestOrder(t *testing.T) {
	bb := New(2)
	x := uint16(0x1122)
	bb.OrderTo(bi.LittleEndian)
	checkIntOrder(t, bb, bi.LittleEndian, x, 0x22, 0x11)

	bb.OrderTo(bi.BigEndian)
	checkIntOrder(t, bb, bi.BigEndian, x, 0x11, 0x22)
}

func checkIntOrder(t *testing.T, bb ByteBuffer, order bi.ByteOrder,
	ui uint16, c, d byte) {
	bb.Clear()
	bb.PutUint16(ui)
	bb.Flip()
	if bb.Order() != order {
		t.Errorf("wanted %v, got %v", order, bb.Order())
	}
	if b := bb.Get(); b != c {
		t.Errorf("wanted %x, got %x", c, b)
	}
	if b := bb.Get(); b != d {
		t.Errorf("wanted %x, got %x", d, b)
	}
}

func TestWrap(t *testing.T) {
	cases := []byte{0, 1, 127, 128, 255}
	bb := Wrap(cases)
	for _, wanted := range cases {
		b := bb.Get()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })
}

func TestByteAccess(t *testing.T) {
	cases := []byte{0, 1, 127, 128, 255}
	bb := New(5)
	bb.OrderTo(bi.LittleEndian)

	for _, c := range cases {
		bb.Put(c)
	}
	checkErrCase(t, ErrOverflow, func() { bb.Put(2) })

	bb.Flip()
	for _, wanted := range cases {
		b := bb.Get()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })
}

func checkErrCase(t *testing.T, expectedErr error, f func()) {
	defer func() {
		if err := recover(); err != expectedErr {
			t.Errorf("ErrCase: wanted %v, got %v", expectedErr, err)
		}
	}()
	f()
}

func TestUint16Access(t *testing.T) {
	cases := []uint16{0, 1, 32767, 32768, math.MaxUint16}
	bb := New(len(cases) * 2)
	bb.OrderTo(bi.LittleEndian)

	for _, c := range cases {
		bb.PutUint16(c)
	}
	checkErrCase(t, ErrOverflow, func() { bb.PutUint16(2) })

	bb.Flip()
	for _, wanted := range cases {
		b := bb.GetUint16()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.GetUint16() })
}

func TestUint32Access(t *testing.T) {
	cases := []uint32{0, 1, 2147483647, 2147483648, math.MaxUint32}
	bb := New(len(cases) * 4)

	for _, c := range cases {
		bb.PutUint32(c)
	}
	checkErrCase(t, ErrOverflow, func() { bb.PutUint32(2) })

	bb.Flip()
	for _, wanted := range cases {
		b := bb.GetUint32()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.GetUint32() })
}

func TestUint64Access(t *testing.T) {
	cases := []uint64{0, 1, 1<<63 - 1, 1 << 63, math.MaxUint64}
	bb := New(len(cases) * 8)

	for _, c := range cases {
		bb.PutUint64(c)
	}
	checkErrCase(t, ErrOverflow, func() { bb.PutUint64(2) })

	bb.Flip()
	for _, wanted := range cases {
		b := bb.GetUint64()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.GetUint64() })
}

func TestCapacity(t *testing.T) {
	bb := New(16)
	checkCursors(t, bb, 16, 0, 16)
}

func TestFlip(t *testing.T) {
	bb := New(16)
	bb.Put(1).PutUint32(0x22).Put(3)
	checkCursors(t, bb, 16, 6, 16)
	bb.Flip()
	checkCursors(t, bb, 16, 0, 6)
}

func TestClear(t *testing.T) {
	bb := New(16)
	bb.Put(1).PutUint32(0x22).Put(3)
	checkCursors(t, bb, 16, 6, 16)
	bb.Clear()
	checkCursors(t, bb, 16, 0, 16)
}

func TestRemaining(t *testing.T) {
	bb := New(16)
	checkCursors(t, bb, 16, 0, 16)

	for bb.HasRemaining() {
		bb.Put(1)
	}
	if bb.Remaining() != 0 {
		t.Errorf("Remaining: wanted %d, got %d", 0, bb.Remaining())
	}
	checkCursors(t, bb, 16, 16, 16)
}

func TestMarkAndReset(t *testing.T) {
	bb := New(16)
	for bb.HasRemaining() {
		bb.Put(1)
	}
	bb.Flip()
	checkErrCase(t, ErrMark, func() { bb.Reset() })

	bb.Mark()
	bb.GetUint32()
	checkCursors(t, bb, 16, 4, 16)
	bb.Reset()
	checkCursors(t, bb, 16, 0, 16)
}

func TestLimit(t *testing.T) {
	bb := New(16)
	for bb.HasRemaining() {
		bb.Put(1)
	}
	bb.Mark()
	checkErrCase(t, ErrLimit, func() { bb.LimitTo(-1) })
	checkErrCase(t, ErrLimit, func() { bb.LimitTo(128) })

	bb.LimitTo(8)
	checkCursors(t, bb, 16, 8, 8)
	checkErrCase(t, ErrMark, func() { bb.Reset() })
}

func TestPosition(t *testing.T) {
	bb := New(16)
	bb.LimitTo(8)
	checkErrCase(t, ErrPosition, func() { bb.PositionTo(-1) })
	checkErrCase(t, ErrPosition, func() { bb.PositionTo(9) })

	bb.PutUint32(0x11).Put(1)
	bb.Mark()
	bb.PositionTo(4)
	checkCursors(t, bb, 16, 4, 8)
	checkErrCase(t, ErrMark, func() { bb.Reset() })
}

func TestRewind(t *testing.T) {
	bb := New(16)
	bb.LimitTo(8)
	bb.PutUint32(0x11).Put(2)
	bb.Mark()
	checkCursors(t, bb, 16, 5, 8)
	bb.Rewind()
	checkCursors(t, bb, 16, 0, 8)
}

func TestCompact(t *testing.T) {
	bb := New(16)
	for bb.HasRemaining() {
		bb.Put(1)
	}
	bb.Flip()
	bb.GetUint32()
	bb.GetUint32()
	checkCursors(t, bb, 16, 8, 16)
	bb.Compact()
	checkCursors(t, bb, 16, 8, 16)
	checkErrCase(t, ErrMark, func() { bb.Reset() })

	bb.Put(2).Put(2)
	bb.Flip()
	for i := 0; i < 8; i++ {
		b := bb.Get()
		if b != 1 {
			t.Errorf("Element: wanted 1, got %d", b)
		}
	}
	for i := 0; i < 2; i++ {
		b := bb.Get()
		if b != 2 {
			t.Errorf("Element: wanted 2, got %d", b)
		}
	}
}

func checkCursors(t *testing.T, bb ByteBuffer, capacity, position, limit int) {
	if bb.Capacity() != capacity {
		t.Errorf("Capacity: wanted: %d, got: %d\n", capacity, bb.Capacity())
	}
	if bb.Position() != position {
		t.Errorf("Position: wanted: %d, got: %d\n", position, bb.Position())
	}
	if bb.Limit() != limit {
		t.Errorf("Limit: wanted: %d, got: %d\n", limit, bb.Limit())
	}
}

func TestByteOrder(t *testing.T) {
	ui := uint16(0x1122)
	if byte(0x22) == (*[2]byte)(us.Pointer(&ui))[0] {
		if ByteOrder != bi.LittleEndian {
			t.Errorf("wanted %v, got %v\n", bi.LittleEndian, ByteOrder)
		}
	} else {
		if ByteOrder != bi.BigEndian {
			t.Errorf("wanted %v, got %v\n", bi.BigEndian, ByteOrder)
		}
	}
}
