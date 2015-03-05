package prime

import (
	"fmt"
	"testing"
)

func BenchmarkAll_10_Eratosthenes(b *testing.B) {
	benchmarkAll(b, SieveEratosthenes, 10)
}

func BenchmarkAll_100_Eratosthenes(b *testing.B) {
	benchmarkAll(b, SieveEratosthenes, 100)
}

func BenchmarkAll_1000_Eratosthenes(b *testing.B) {
	benchmarkAll(b, SieveEratosthenes, 1000)
}

func BenchmarkAll_10000_Eratosthenes(b *testing.B) {
	benchmarkAll(b, SieveEratosthenes, 10000)
}

func BenchmarkAll_10_EratosthenesO(b *testing.B) {
	benchmarkAll(b, SieveEratosthenesO, 10)
}

func BenchmarkAll_100_EratosthenesO(b *testing.B) {
	benchmarkAll(b, SieveEratosthenesO, 100)
}

func BenchmarkAll_1000_EratosthenesO(b *testing.B) {
	benchmarkAll(b, SieveEratosthenesO, 1000)
}

func BenchmarkAll_10000_EratosthenesO(b *testing.B) {
	benchmarkAll(b, SieveEratosthenesO, 10000)
}

func benchmarkAll(b *testing.B, ps Sieve, upTo uint64) {
	for i := 0; i < b.N; i++ {
		ch := ps.All(upTo)
		for _ = range ch {
		}
	}
	b.ReportAllocs()
}

func TestAll_10_Eratosthenes(t *testing.T) {
	testAll(t, SieveEratosthenes, 10)
}

func TestAll_10_EratosthenesO(t *testing.T) {
	testAll(t, SieveEratosthenesO, 10)
}

func TestAll_10_Atkin(t *testing.T) {
	testAll(t, SieveAtkin, 10)
}

func testAll(t *testing.T, ps Sieve, upTo uint64) {
	ch := ps.All(upTo)
	s := []uint64(nil)
	for p := range ch {
		s = append(s, p)
	}

	c := []uint64{2, 3, 5, 7}
	if fmt.Sprintf("%v", s) != fmt.Sprintf("%v", c) {
		t.Errorf("wanted: %v, got: %v\n", c, s)
	}
}
