package pair

import (
	"reflect"
	"sort"
	"testing"
)

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
