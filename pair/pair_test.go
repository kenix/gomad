package pair

import (
	"reflect"
	"sort"
	"testing"
)

func TestClosestPair(t *testing.T) {
	cases := []struct {
		in  []*Point
		out float64
	}{
		{
			[]*Point{
				NewPoint(2, 3), NewPoint(12, 30), NewPoint(40, 50),
				NewPoint(5, 1), NewPoint(12, 10), NewPoint(3, 4),
				NewPoint(50, 10), NewPoint(120, 100), NewPoint(30, 40),
				NewPoint(52, 12), NewPoint(122, 102), NewPoint(32, 42),
				NewPoint(54, 14), NewPoint(124, 104), NewPoint(34, 44),
				NewPoint(56, 16), NewPoint(126, 106), NewPoint(36, 46),
				NewPoint(58, 18), NewPoint(128, 108), NewPoint(38, 48),
				NewPoint(60, 20), NewPoint(130, 110),
			},
			1.4142135623730951,
		},
	}

	for _, c := range cases {
		p1, p2, d := ClosestPair(c.in)
		if d != c.out {
			t.Errorf("wanted %v got %v %s %s\n", c.out, d, p1, p2)
		}
	}
}

func TestIndices(t *testing.T) {
	cases := []struct {
		in  Indices
		out []int
	}{
		{
			NewIndices(ByX{NewPoints([]*Point{NewPoint(2, 1), NewPoint(0, 3), NewPoint(1, 2)})}),
			[]int{1, 2, 0},
		},
		{
			NewIndices(ByY{NewPoints([]*Point{NewPoint(2, 1), NewPoint(0, 3), NewPoint(1, 2)})}),
			[]int{0, 2, 1},
		},
	}

	for _, c := range cases {
		sort.Sort(c.in)
		if !reflect.DeepEqual(c.in.Indices, c.out) {
			t.Errorf("wanted %v got %v\n", c.out, c.in.Indices)
		}
	}
}

func TestPointsSort(t *testing.T) {
	cases := []struct {
		in   Points
		outX Points
		outY Points
	}{
		{
			NewPoints([]*Point{NewPoint(2, 1), NewPoint(0, 3), NewPoint(1, 2)}),
			NewPoints([]*Point{NewPoint(0, 3), NewPoint(1, 2), NewPoint(2, 1)}),
			NewPoints([]*Point{NewPoint(2, 1), NewPoint(1, 2), NewPoint(0, 3)}),
		},
	}

	for _, c := range cases {
		sort.Sort(ByX{c.in})
		if !reflect.DeepEqual(c.in, c.outX) {
			t.Errorf("byX wanted %v got %v\n", c.outX, c.in)
		}
		sort.Sort(ByY{c.in})
		if !reflect.DeepEqual(c.in, c.outY) {
			t.Errorf("byY wanted %v got %v\n", c.outY, c.in)
		}
	}
}

func TestDistance(t *testing.T) {
	cases := []struct {
		p1, p2 *Point
		dist   float64
	}{
		{NewPoint(0.3, 0), NewPoint(0, 0.4), 0.5},
		{NewPoint(-3, 0), NewPoint(0, 4), 5},
		{NewPoint(-0.3, 0), NewPoint(0, -0.4), 0.5},
		{NewPoint(3, 0), NewPoint(0, -4), 5},
	}

	for _, c := range cases {
		assertEqual(t, c.p1, c.p2, c.dist, c.p1.Distance(c.p2))
		assertEqual(t, c.p2, c.p1, c.dist, c.p2.Distance(c.p1))
		assertEqual(t, c.p2, c.p1, c.dist, Distance(c.p1, c.p2))
	}
}

func assertEqual(t *testing.T, p1, p2 *Point, wanted, got float64) {
	if got != wanted {
		t.Errorf("distance between %s and %s: wanted %f, got %f\n",
			p1, p2, wanted, got)
	}
}
