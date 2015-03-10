package pair

import "math"
import "fmt"

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
	return len([]*Point(ps))
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

func ClosestPair(ps Points) (*Point, *Point) {
	// TODO divide and conquer
	return nil, nil
}
