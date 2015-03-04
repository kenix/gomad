package bytebuffer

import (
	bi "encoding/binary"
	"math"
	"testing"
)

func BenchmarkGet512_C(b *testing.B) {
	benchmarkGet(b, NewBB(1<<20), makeSlice(1<<9))
}

func BenchmarkGet1K_C(b *testing.B) {
	benchmarkGet(b, NewBB(1<<20), makeSlice(1<<10))
}

func BenchmarkGet32K_C(b *testing.B) {
	benchmarkGet(b, NewBB(1<<20), makeSlice(1<<15))
}

func BenchmarkRead512_C(b *testing.B) {
	benchmarkRead(b, NewBB(1<<20), 1<<9)
}

func BenchmarkRead1K_C(b *testing.B) {
	benchmarkRead(b, NewBB(1<<20), 1<<10)
}

func BenchmarkRead32K_C(b *testing.B) {
	benchmarkRead(b, NewBB(1<<20), 1<<15)
}

func BenchmarkPut1K_C(b *testing.B) {
	benchmarkPut(b, NewBB(1<<20), makeSlice(1024))
}

func BenchmarkPut_C(b *testing.B) {
	benchmarkPut(b, NewBB(1<<20), byte(math.MaxUint8))
}

func BenchmarkGet_C(b *testing.B) {
	benchmarkGet(b, NewBB(1<<20), byte(math.MaxUint8))
}

func BenchmarkPutUint16_C(b *testing.B) {
	benchmarkPut(b, NewBB(1<<20), uint16(math.MaxUint16))
}

func BenchmarkGetUint16_C(b *testing.B) {
	benchmarkGet(b, NewBB(1<<20), uint16(math.MaxUint16))
}

func TestWrap_C(t *testing.T) {
	cases := []byte{0, 1, 127, 128, 255}
	bb := WrapBB(cases)
	for _, wanted := range cases {
		b := bb.Get()
		if wanted != b {
			t.Errorf("wanted:%d, got:%d\n", wanted, b)
		}
	}
	checkErrCase(t, ErrUnderflow, func() { bb.Get() })
}

func TestByteAccess_C(t *testing.T) {
	cases := []byte{0, 1, 127, 128, 255}
	bb := NewBB(5)
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

func TestUint16Access_C(t *testing.T) {
	cases := []uint16{0, 1, 32767, 32768, math.MaxUint16}
	bb := NewBB(len(cases) * 2)
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
