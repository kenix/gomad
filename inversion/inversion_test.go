package inversion

import (
	"github.com/mentopolis/gomad/util"
	"testing"
)

func BenchmarkInv_10K_DivCon(b *testing.B) {
	benchmarkInv(b, InvDivCon, shuffle(10000))
}

func BenchmarkInv_100K_DivCon(b *testing.B) {
	benchmarkInv(b, InvDivCon, shuffle(100000))
}

func shuffle(n int) []int64 {
	s := make([]int64, n, n)
	for i := 0; i < n; i++ {
		s[i] = int64(i)
	}
	util.Shuffle(s)
	return s
}

func benchmarkInv(b *testing.B, inv Inversion, a []int64) {
	for i := 0; i < b.N; i++ {
		inv.Count(a)
	}
}

func TestInvDivCon(t *testing.T) {
	cases := []struct {
		in  []int64
		out uint64
	}{
		{[]int64{1, 3, 5, 2, 4, 6}, 3},
		{[]int64{9, 6, 3, 5, 4, 1, 0, 7, 8, 2}, 28},
	}

	inv := InvDivCon
	for _, c := range cases {
		got := inv.Count(c.in)
		if c.out != got {
			t.Errorf("wanted %d, got %d\n", c.out, got)
		}
	}
}
