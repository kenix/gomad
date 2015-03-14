package pair

import (
	"fmt"
	"math"
	"sort"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (p *Point) String() string {
	return fmt.Sprintf("P(%f,%f)", p.x, p.y)
}

func (p *Point) X() float64 {
	return p.x
}

func (p *Point) Y() float64 {
	return p.y
}

func (p *Point) Distance(p2 *Point) float64 {
	return math.Sqrt(math.Pow(p.x-p2.x, 2) + math.Pow(p.y-p2.y, 2))
}

func Distance(p1, p2 *Point) float64 {
	return p1.Distance(p2)
}

type Points []*Point

func NewPoints(ps []*Point) Points {
	return Points(ps)
}

func (ps Points) Len() int {
	return len(ps)
}

func (ps Points) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

type ByX struct {
	Points
}

func (ps ByX) Less(i, j int) bool {
	return ps.Points[i].x < ps.Points[j].x
}

type ByY struct {
	Points
}

func (ps ByY) Less(i, j int) bool {
	return ps.Points[i].y < ps.Points[j].y
}

type Indices struct {
	Indices []int
	sort.Interface
}

func NewIndices(underlying sort.Interface) Indices {
	indices := []int(nil)
	for i := 0; i < underlying.Len(); i++ {
		indices = append(indices, i)
	}
	return Indices{indices, underlying}
}

func (ids Indices) Swap(i, j int) {
	ids.Indices[i], ids.Indices[j] = ids.Indices[j], ids.Indices[i]
}

func (ids Indices) Less(i, j int) bool {
	return ids.Interface.Less(ids.Indices[i], ids.Indices[j])
}

func ClosestPair(ps []*Point) (*Point, *Point, float64) {
	switch len(ps) {
	case 0, 1:
		return nil, nil, math.MaxFloat64
	case 2:
		return ps[0], ps[1], Distance(ps[0], ps[1])
	}

	pa := NewPoints(ps)
	sort.Sort(ByX{pa})
	pyi := NewIndices(ByY{pa})
	sort.Sort(pyi)

	return dccPair(ps, pyi.Indices, pa)
}

func dccPair(ps []*Point, pyi []int, pa Points) (*Point, *Point, float64) {
	switch len(ps) {
	case 2:
		return ps[0], ps[1], Distance(ps[0], ps[1])
	case 3:
		d01 := Distance(ps[0], ps[1])
		d02 := Distance(ps[0], ps[2])
		d12 := Distance(ps[1], ps[2])
		if d01 <= d02 {
			if d12 <= d01 {
				return ps[1], ps[2], d12
			}
			return ps[0], ps[1], d01
		}
		if d12 <= d02 {
			return ps[1], ps[2], d12
		}
		return ps[0], ps[2], d02
	}

	// TODO use goroutine
	pl1, pl2, dl := dccPair(ps[:len(ps)/2], pyi, pa)
	pr1, pr2, dr := dccPair(ps[len(ps)/2:], pyi, pa)
	d := math.Min(dl, dr)

	sp := []*Point(nil)
	px := ps[len(ps)/2]
	for i := 0; i < len(pyi); i++ {
		p := pa[pyi[i]]
		if math.Abs(p.x-px.x) <= d && math.Abs(p.y-px.y) <= d {
			sp = append(sp, p)
		}
	}
	ps1, ps2, ds := closestSplitPair(sp)
	if ds <= d {
		return ps1, ps2, ds
	}
	if dl <= dr {
		return pl1, pl2, dl
	}
	return pr1, pr2, dr
}

func closestSplitPair(ps []*Point) (ps1 *Point, ps2 *Point, d float64) {
	switch len(ps) {
	case 0, 1:
		return nil, nil, math.MaxFloat64
	case 2:
		return ps[0], ps[1], Distance(ps[0], ps[1])
	}

	d = math.MaxFloat64
	n := int(math.Min(7, float64(len(ps))))
	for i := 0; i <= len(ps)-n; i++ {
		for j := i + 1; j < n; j++ {
			ds := Distance(ps[i], ps[j])
			if ds < d {
				ps1, ps2, d = ps[i], ps[j], ds
			}
		}
	}
	return
}
