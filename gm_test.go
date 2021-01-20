package gm

import (
	"math"
	"testing"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

const (
	pi    = math.Pi
	sqrt2 = math.Sqrt2
	sqrt3 = 1.73205080756887729352744634150587236694280525381038062805580698 // https://oeis.org/A002194
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		p, n s2.LatLng
		gm   *GeneralizedMercator
	}{
		// Antipodes on the great circle of longitude ±90° (X == 0)
		{
			p: s2.LatLng{Lat: pi / 2, Lng: 0},
			n: s2.LatLng{Lat: -pi / 2, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, 1},
				neg: r3.Vector{0, 0, -1},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{0, 0, 1},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: -pi / 2, Lng: 0},
			n: s2.LatLng{Lat: pi / 2, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, -1},
				neg: r3.Vector{0, 0, 1},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{0, 0, -1},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: pi / 2},
			n: s2.LatLng{Lat: 0, Lng: -pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 1, 0},
				neg: r3.Vector{0, -1, 0},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 0, -1},
				k:   r3.Vector{0, 1, 0},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: -pi / 2},
			n: s2.LatLng{Lat: 0, Lng: pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, -1, 0},
				neg: r3.Vector{0, 1, 0},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 0, 1},
				k:   r3.Vector{0, -1, 0},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: pi / 4, Lng: -pi / 2},
			n: s2.LatLng{Lat: -pi / 4, Lng: pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, -sqrt2 / 2, sqrt2 / 2},
				neg: r3.Vector{0, sqrt2 / 2, -sqrt2 / 2},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, sqrt2 / 2, sqrt2 / 2},
				k:   r3.Vector{0, -sqrt2 / 2, sqrt2 / 2},
				d:   math.Inf(1),
			},
		},
		// Antipodes elsewhere on the Equator (X != 0, Z == 0)
		{
			p: s2.LatLng{Lat: 0, Lng: 0},
			n: s2.LatLng{Lat: 0, Lng: pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{1, 0, 0},
				neg: r3.Vector{-1, 0, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{1, 0, 0},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: pi},
			n: s2.LatLng{Lat: 0, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{-1, 0, 0},
				neg: r3.Vector{1, 0, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{-1, 0, 0},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: pi / 4},
			n: s2.LatLng{Lat: 0, Lng: -3 * pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{sqrt2 / 2, sqrt2 / 2, 0},
				neg: r3.Vector{-sqrt2 / 2, -sqrt2 / 2, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{sqrt2 / 2, -sqrt2 / 2, 0},
				k:   r3.Vector{sqrt2 / 2, sqrt2 / 2, 0},
				d:   math.Inf(1),
			},
		},
		// Other antipodes (X != 0, Z != 0)
		{
			p: s2.LatLng{Lat: pi / 4, Lng: 0},
			n: s2.LatLng{Lat: -pi / 4, Lng: pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{sqrt2 / 2, 0, sqrt2 / 2},
				neg: r3.Vector{-sqrt2 / 2, 0, -sqrt2 / 2},
				i:   r3.Vector{sqrt2 / 2, 0, -sqrt2 / 2},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{sqrt2 / 2, 0, sqrt2 / 2},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: -pi / 4, Lng: 0},
			n: s2.LatLng{Lat: pi / 4, Lng: pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{sqrt2 / 2, 0, -sqrt2 / 2},
				neg: r3.Vector{-sqrt2 / 2, 0, sqrt2 / 2},
				i:   r3.Vector{sqrt2 / 2, 0, sqrt2 / 2},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{sqrt2 / 2, 0, -sqrt2 / 2},
				d:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: pi / 4, Lng: pi / 4},
			n: s2.LatLng{Lat: -pi / 4, Lng: -3 * pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0.5, 0.5, sqrt2 / 2},
				neg: r3.Vector{-0.5, -0.5, -sqrt2 / 2},
				i:   r3.Vector{sqrt2 / sqrt3, 0, -1 / sqrt3},
				j:   r3.Vector{-1 / (2 * sqrt3), sqrt3 / 2, -1 / (sqrt2 * sqrt3)},
				k:   r3.Vector{0.5, 0.5, sqrt2 / 2},
				d:   math.Inf(1),
			},
		},
		// Non-antipodes
		{
			p: s2.LatLng{Lat: pi / 3, Lng: 0},
			n: s2.LatLng{Lat: -pi / 3, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0.5, 0, sqrt3 / 2},
				neg: r3.Vector{0.5, 0, -sqrt3 / 2},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{0, 0, 1},
				d:   2,
			},
		},
		{
			p: s2.LatLng{Lat: pi / 3, Lng: pi / 4},
			n: s2.LatLng{Lat: -pi / 3, Lng: pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{sqrt2 / 4, sqrt2 / 4, sqrt3 / 2},
				neg: r3.Vector{sqrt2 / 4, sqrt2 / 4, -sqrt3 / 2},
				i:   r3.Vector{sqrt2 / 2, sqrt2 / 2, 0},
				j:   r3.Vector{-sqrt2 / 2, sqrt2 / 2, 0},
				k:   r3.Vector{0, 0, 1},
				d:   2,
			},
		},
		{
			p: s2.LatLng{Lat: pi / 2, Lng: 0},
			n: s2.LatLng{Lat: -pi / 6, Lng: 2 * pi / 3},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, 1},
				neg: r3.Vector{-sqrt3 / 4, 0.75, -0.5},
				i:   r3.Vector{-sqrt3 / 4, 0.75, 0.5},
				j:   r3.Vector{-sqrt3 / 2, -0.5, 0},
				k:   r3.Vector{0.25, -sqrt3 / 4, sqrt3 / 2},
				d:   2,
			},
		},
		{
			p: s2.LatLng{Lat: -pi / 4, Lng: 0},
			n: s2.LatLng{Lat: -pi / 4, Lng: pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{sqrt2 / 2, 0, -sqrt2 / 2},
				neg: r3.Vector{-sqrt2 / 2, 0, -sqrt2 / 2},
				i:   r3.Vector{0, 0, -1},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{1, 0, 0},
				d:   sqrt2,
			},
		},
	} {
		if gm := New(test.p, test.n); !gmApproxEqual(gm, test.gm) {
			t.Errorf("New(%v, %v): got %+v, want %+v", test.p, test.n, gm, test.gm)
		}
	}
}

type proj struct {
	s s2.LatLng
	r r2.Point
}

var projTests = []struct {
	gm *GeneralizedMercator
	ps []proj
}{
	{
		gm: &GeneralizedMercator{
			pos: r3.Vector{0, 0, 1},
			neg: r3.Vector{0, 0, -1},
			i:   r3.Vector{1, 0, 0},
			j:   r3.Vector{0, 1, 0},
			k:   r3.Vector{0, 0, 1},
			d:   math.Inf(1),
		},
		ps: []proj{
			{s2.LatLng{Lat: pi / 2}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, math.Inf(1)).X, Y: p0(0.5, math.Inf(1)).Y, Z: p0(0.5, math.Inf(1)).Z}}), r2.Point{X: 0, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, math.Inf(1)).X, Y: p1(0.5, math.Inf(1)).Y, Z: p1(0.5, math.Inf(1)).Z}}), r2.Point{X: pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, math.Inf(1)).X, Y: p2(0.5, math.Inf(1)).Y, Z: p2(0.5, math.Inf(1)).Z}}), r2.Point{X: pi, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, math.Inf(1)).X, Y: p3(0.5, math.Inf(1)).Y, Z: p3(0.5, math.Inf(1)).Z}}), r2.Point{X: -pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLng{Lat: 0, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: pi / 2}, r2.Point{X: pi / 2, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: pi}, r2.Point{X: pi, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: -pi / 2}, r2.Point{X: -pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, math.Inf(1)).X, Y: p0(-0.5, math.Inf(1)).Y, Z: p0(-0.5, math.Inf(1)).Z}}), r2.Point{X: 0, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, math.Inf(1)).X, Y: p1(-0.5, math.Inf(1)).Y, Z: p1(-0.5, math.Inf(1)).Z}}), r2.Point{X: pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, math.Inf(1)).X, Y: p2(-0.5, math.Inf(1)).Y, Z: p2(-0.5, math.Inf(1)).Z}}), r2.Point{X: pi, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, math.Inf(1)).X, Y: p3(-0.5, math.Inf(1)).Y, Z: p3(-0.5, math.Inf(1)).Z}}), r2.Point{X: -pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLng{Lat: -pi / 2}, r2.Point{Y: math.Inf(-1)}},
		},
	},
	{
		gm: &GeneralizedMercator{
			pos: r3.Vector{0.5, -sqrt3 / 2, 0},
			neg: r3.Vector{0.5, sqrt3 / 2, 0},
			i:   r3.Vector{1, 0, 0},
			j:   r3.Vector{0, 0, 1},
			k:   r3.Vector{0, -1, 0},
			d:   2,
		},
		ps: []proj{
			{s2.LatLng{Lat: 0, Lng: -pi / 3}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, 2).X, Y: -p0(0.5, 2).Z, Z: p0(0.5, 2).Y}}), r2.Point{X: 0, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, 2).X, Y: -p1(0.5, 2).Z, Z: p1(0.5, 2).Y}}), r2.Point{X: pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, 2).X, Y: -p2(0.5, 2).Z, Z: p2(0.5, 2).Y}}), r2.Point{X: pi, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, 2).X, Y: -p3(0.5, 2).Z, Z: p3(0.5, 2).Y}}), r2.Point{X: -pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLng{Lat: 0, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: pi / 2}, r2.Point{X: pi / 2, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: pi}, r2.Point{X: pi, Y: 0}},
			{s2.LatLng{Lat: -pi / 2}, r2.Point{X: -pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, 2).X, Y: -p0(-0.5, 2).Z, Z: p0(-0.5, 2).Y}}), r2.Point{X: 0, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, 2).X, Y: -p1(-0.5, 2).Z, Z: p1(-0.5, 2).Y}}), r2.Point{X: pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, 2).X, Y: -p2(-0.5, 2).Z, Z: p2(-0.5, 2).Y}}), r2.Point{X: pi, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, 2).X, Y: -p3(-0.5, 2).Z, Z: p3(-0.5, 2).Y}}), r2.Point{X: -pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLng{Lat: 0, Lng: pi / 3}, r2.Point{Y: math.Inf(-1)}},
		},
	},
	{
		gm: &GeneralizedMercator{
			pos: r3.Vector{sqrt2 / 2, 0, -sqrt2 / 2},
			neg: r3.Vector{-sqrt2 / 2, 0, -sqrt2 / 2},
			i:   r3.Vector{0, 0, -1},
			j:   r3.Vector{0, 1, 0},
			k:   r3.Vector{1, 0, 0},
			d:   sqrt2,
		},
		ps: []proj{
			{s2.LatLng{Lat: -pi / 4, Lng: 0}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, sqrt2).Z, Y: p0(0.5, sqrt2).Y, Z: -p0(0.5, sqrt2).X}}), r2.Point{X: 0, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, sqrt2).Z, Y: p1(0.5, sqrt2).Y, Z: -p1(0.5, sqrt2).X}}), r2.Point{X: pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, sqrt2).Z, Y: p2(0.5, sqrt2).Y, Z: -p2(0.5, sqrt2).X}}), r2.Point{X: pi, Y: math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, sqrt2).Z, Y: p3(0.5, sqrt2).Y, Z: -p3(0.5, sqrt2).X}}), r2.Point{X: -pi / 2, Y: math.Log(sqrt3)}},
			{s2.LatLng{Lat: -pi / 2, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: pi / 2}, r2.Point{X: pi / 2, Y: 0}},
			{s2.LatLng{Lat: pi / 2, Lng: 0}, r2.Point{X: pi, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: -pi / 2}, r2.Point{X: -pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, sqrt2).Z, Y: p0(-0.5, sqrt2).Y, Z: -p0(-0.5, sqrt2).X}}), r2.Point{X: 0, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, sqrt2).Z, Y: p1(-0.5, sqrt2).Y, Z: -p1(-0.5, sqrt2).X}}), r2.Point{X: pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, sqrt2).Z, Y: p2(-0.5, sqrt2).Y, Z: -p2(-0.5, sqrt2).X}}), r2.Point{X: pi, Y: -math.Log(sqrt3)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, sqrt2).Z, Y: p3(-0.5, sqrt2).Y, Z: -p3(-0.5, sqrt2).X}}), r2.Point{X: -pi / 2, Y: -math.Log(sqrt3)}},
			{s2.LatLng{Lat: -pi / 4, Lng: pi}, r2.Point{Y: math.Inf(-1)}},
		},
	},
}

func TestProject(t *testing.T) {
	for _, test := range projTests {
		for _, p := range test.ps {
			if got := test.gm.Project(p.s); !ptApproxEqual(got, p.r) {
				t.Errorf("Project(%+v, %+v): got %+v, want %+v", test.gm, p.s, got, p.r)
			}
		}
	}
}

func TestUnproject(t *testing.T) {
	for _, test := range projTests {
		for _, p := range test.ps {
			if got := test.gm.Unproject(p.r); !llApproxEqual(got, p.s) {
				t.Errorf("Unproject(%+v, %+v): got %+v, want %+v", test.gm, p.r, got, p.s)
			}
		}
	}
}

func ptApproxEqual(a, b r2.Point) bool {
	return (a.X == b.X || math.Abs(a.X-b.X) < 1e-15) && (a.Y == b.Y || math.Abs(a.Y-b.Y) < 1e-15)
}

func llApproxEqual(a, b s2.LatLng) bool {
	return approxEqual(s2.PointFromLatLng(a).Vector, s2.PointFromLatLng(b).Vector)
}

func gmApproxEqual(a, b *GeneralizedMercator) bool {
	return approxEqual(a.pos, b.pos) &&
		approxEqual(a.neg, b.neg) &&
		approxEqual(a.i, b.i) &&
		approxEqual(a.j, b.j) &&
		approxEqual(a.k, b.k) &&
		(a.d == b.d || math.Abs(a.d-b.d) < 1e-15)
}

// p0, p1, p2, and p3 describe points on the circular intersection of the unit sphere
// with the plane containing x == d that lies at a signed distance of c from the origin:
//  point	extremum	projective longitude
// 	p0		max x		0
// 	p1		max y		pi/2
// 	p2		min x		pi
// 	p3		min y		-pi/2

func p0(c, d float64) r3.Vector {
	var psi, beta = math.Asin(c), math.Asin(c / d)
	return r3.Vector{
		X: math.Sin(psi)*math.Sin(beta) + math.Cos(psi)*math.Cos(beta),
		Y: 0,
		Z: math.Sin(psi)*math.Cos(beta) - math.Cos(psi)*math.Sin(beta),
	}
}

func p1(c, d float64) r3.Vector {
	var psi, beta = math.Asin(c), math.Asin(c / d)
	return r3.Vector{
		X: math.Sin(psi) * math.Sin(beta),
		Y: math.Cos(psi),
		Z: math.Sin(psi) * math.Cos(beta),
	}
}

func p2(c, d float64) r3.Vector {
	var psi, beta = math.Asin(c), math.Asin(c / d)
	return r3.Vector{
		X: math.Sin(psi)*math.Sin(beta) - math.Cos(psi)*math.Cos(beta),
		Y: 0,
		Z: math.Sin(psi)*math.Cos(beta) + math.Cos(psi)*math.Sin(beta),
	}
}

func p3(c, d float64) r3.Vector {
	var psi, beta = math.Asin(c), math.Asin(c / d)
	return r3.Vector{
		X: math.Sin(psi) * math.Sin(beta),
		Y: -math.Cos(psi),
		Z: math.Sin(psi) * math.Cos(beta),
	}
}

func BenchmarkProject(b *testing.B) {
	gm := New(s2.LatLng{Lat: math.Pi / 3}, s2.LatLng{Lat: -math.Pi / 3})
	for i := 0; i < b.N; i++ {
		gm.Project(s2.LatLng{Lat: math.Pi / 4, Lng: 3 * math.Pi / 4})
	}
}

func BenchmarkUnproject(b *testing.B) {
	gm := New(s2.LatLng{Lat: math.Pi / 3}, s2.LatLng{Lat: -math.Pi / 3})
	for i := 0; i < b.N; i++ {
		gm.Unproject(r2.Point{1, 1})
	}
}
