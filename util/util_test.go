package util

import (
	"math"
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	cases := []struct {
		in  []int
		out []int
	}{
		{
			[]int{3, 9, 8, 2, 5, 7, 1, 6, 4},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			[]int{3, 8, 2, 5, 9, 7, 1, 6, 4},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			[]int{9, 8, 7, 6, 5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, c := range cases {
		QuickSort(c.in)
		if !reflect.DeepEqual(c.in, c.out) {
			t.Errorf("wanted %v, got %v\n", c.out, c.in)
		}
	}
}

func BenchmarkLeadingZeros(b *testing.B) {
	s := []interface{}{byte(1), int8(1), uint16(1), int16(1), uint16(1),
		rune(1), int32(1), uint32(1), int64(1), uint64(1)}
	for i := 0; i < b.N; i++ {
		for _, c := range s {
			if LeadingZeros(c) < 0 {
				b.Fail()
			}
		}
	}
}

func BenchmarkSize(b *testing.B) {
	s := []interface{}{byte(1), int8(1), uint16(1), int16(1), uint16(1),
		rune(1), int32(1), uint32(1), int64(1), uint64(1)}
	for i := 0; i < b.N; i++ {
		for _, c := range s {
			if Size(c) < 0 {
				b.Fail()
			}
		}
	}
}

func TestLeadingZeros(t *testing.T) {
	cases := []struct {
		in  interface{}
		out int
	}{
		{byte(0x00), 8},
		{byte(0x08), 4},
		{byte(math.MaxUint8), 0},
		{uint8(0x00), 8},
		{uint8(0x08), 4},
		{uint8(math.MaxUint8), 0},
		{int8(0x00), 8},
		{int8(0x08), 4},
		{int8(math.MinInt8), 0},
		{uint16(0x0000), 16},
		{uint16(0x0080), 8},
		{uint16(math.MaxUint16), 0},
		{int16(0x0000), 16},
		{int16(0x0080), 8},
		{int16(math.MinInt16), 0},
		{rune(0x00000000), 32},
		{rune(0x00008000), 16},
		{rune(math.MinInt32), 0},
		{uint32(0x00000000), 32},
		{uint32(0x00008000), 16},
		{uint32(math.MaxUint32), 0},
		{int32(0x00000000), 32},
		{int32(0x00008000), 16},
		{int32(math.MinInt32), 0},
		{uint64(0x0000000000000000), 64},
		{uint64(0x0000000080000000), 32},
		{uint64(math.MaxUint64), 0},
		{int64(0x0000000000000000), 64},
		{int64(0x0000000080000000), 32},
		{int64(math.MinInt64), 0},
		{int(1), -1},
		{uint(1), -1},
		{"blah", -1},
	}
	for _, c := range cases {
		got := LeadingZeros(c.in)
		if got != c.out {
			t.Errorf("LeadingZeros(%x): wanted %d, got %d", c.in, c.out, got)
		}
	}

}

func TestSize(t *testing.T) {
	cases := []struct {
		val  interface{}
		size int
	}{
		{byte(5), 8},
		{uint8(5), 8},
		{int8(-5), 8},
		{uint16(5), 16},
		{int16(-5), 16},
		{int32(-5), 32},
		{uint32(5), 32},
		{rune(5), 32},
		{int64(5), 64},
		{uint64(5), 64},
		{int(5), -1},
		{uint(5), -1},
	}
	for _, c := range cases {
		s := Size(c.val)
		if s != c.size {
			t.Errorf("%t size %d, got %d", c.val, c.size, s)
		}
	}
}
