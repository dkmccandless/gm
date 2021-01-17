package gm

import (
	"math"
	"testing"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		p, n s2.LatLng
		gm   *GeneralizedMercator
	}{
		// Antipodes on the great circle of longitude ±90° (X == 0)
		{
			p: s2.LatLng{Lat: math.Pi / 2, Lng: 0},
			n: s2.LatLng{Lat: -math.Pi / 2, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, 1},
				neg: r3.Vector{0, 0, -1},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{0, 0, 1},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: -math.Pi / 2, Lng: 0},
			n: s2.LatLng{Lat: math.Pi / 2, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, -1},
				neg: r3.Vector{0, 0, 1},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{0, 0, -1},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: math.Pi / 2},
			n: s2.LatLng{Lat: 0, Lng: -math.Pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 1, 0},
				neg: r3.Vector{0, -1, 0},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 0, -1},
				k:   r3.Vector{0, 1, 0},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: -math.Pi / 2},
			n: s2.LatLng{Lat: 0, Lng: math.Pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, -1, 0},
				neg: r3.Vector{0, 1, 0},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 0, 1},
				k:   r3.Vector{0, -1, 0},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: math.Pi / 4, Lng: -math.Pi / 2},
			n: s2.LatLng{Lat: -math.Pi / 4, Lng: math.Pi / 2},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, -math.Sqrt2 / 2, math.Sqrt2 / 2},
				neg: r3.Vector{0, math.Sqrt2 / 2, -math.Sqrt2 / 2},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, math.Sqrt2 / 2, math.Sqrt2 / 2},
				k:   r3.Vector{0, -math.Sqrt2 / 2, math.Sqrt2 / 2},
				t:   math.Inf(1),
			},
		},
		// Antipodes elsewhere on the Equator (X != 0, Z == 0)
		{
			p: s2.LatLng{Lat: 0, Lng: 0},
			n: s2.LatLng{Lat: 0, Lng: math.Pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{1, 0, 0},
				neg: r3.Vector{-1, 0, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{1, 0, 0},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: math.Pi},
			n: s2.LatLng{Lat: 0, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{-1, 0, 0},
				neg: r3.Vector{1, 0, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{-1, 0, 0},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: 0, Lng: math.Pi / 4},
			n: s2.LatLng{Lat: 0, Lng: -3 * math.Pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{math.Sqrt2 / 2, math.Sqrt2 / 2, 0},
				neg: r3.Vector{-math.Sqrt2 / 2, -math.Sqrt2 / 2, 0},
				i:   r3.Vector{0, 0, 1},
				j:   r3.Vector{math.Sqrt2 / 2, -math.Sqrt2 / 2, 0},
				k:   r3.Vector{math.Sqrt2 / 2, math.Sqrt2 / 2, 0},
				t:   math.Inf(1),
			},
		},
		// Other antipodes (X != 0, Z != 0)
		{
			p: s2.LatLng{Lat: math.Pi / 4, Lng: 0},
			n: s2.LatLng{Lat: -math.Pi / 4, Lng: math.Pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{math.Sqrt2 / 2, 0, math.Sqrt2 / 2},
				neg: r3.Vector{-math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				i:   r3.Vector{math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{math.Sqrt2 / 2, 0, math.Sqrt2 / 2},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: -math.Pi / 4, Lng: 0},
			n: s2.LatLng{Lat: math.Pi / 4, Lng: math.Pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				neg: r3.Vector{-math.Sqrt2 / 2, 0, math.Sqrt2 / 2},
				i:   r3.Vector{math.Sqrt2 / 2, 0, math.Sqrt2 / 2},
				j:   r3.Vector{0, -1, 0},
				k:   r3.Vector{math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				t:   math.Inf(1),
			},
		},
		{
			p: s2.LatLng{Lat: math.Pi / 4, Lng: math.Pi / 4},
			n: s2.LatLng{Lat: -math.Pi / 4, Lng: -3 * math.Pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0.5, 0.5, math.Sqrt2 / 2},
				neg: r3.Vector{-0.5, -0.5, -math.Sqrt2 / 2},
				i:   r3.Vector{math.Sqrt(2 / 3.), 0, -math.Sqrt(1 / 3.)},
				j:   r3.Vector{-1 / (2 * math.Sqrt(3)), math.Sqrt(3) / 2, -1 / math.Sqrt(6)},
				k:   r3.Vector{0.5, 0.5, math.Sqrt2 / 2},
				t:   math.Inf(1),
			},
		},
		// Non-antipodes
		{
			p: s2.LatLng{Lat: math.Pi / 3, Lng: 0},
			n: s2.LatLng{Lat: -math.Pi / 3, Lng: 0},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0.5, 0, math.Sqrt(3) / 2},
				neg: r3.Vector{0.5, 0, -math.Sqrt(3) / 2},
				i:   r3.Vector{1, 0, 0},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{0, 0, 1},
				t:   2,
			},
		},
		{
			p: s2.LatLng{Lat: math.Pi / 3, Lng: math.Pi / 4},
			n: s2.LatLng{Lat: -math.Pi / 3, Lng: math.Pi / 4},
			gm: &GeneralizedMercator{
				pos: r3.Vector{math.Sqrt2 / 4, math.Sqrt2 / 4, math.Sqrt(3) / 2},
				neg: r3.Vector{math.Sqrt2 / 4, math.Sqrt2 / 4, -math.Sqrt(3) / 2},
				i:   r3.Vector{math.Sqrt2 / 2, math.Sqrt2 / 2, 0},
				j:   r3.Vector{-math.Sqrt2 / 2, math.Sqrt2 / 2, 0},
				k:   r3.Vector{0, 0, 1},
				t:   2,
			},
		},
		{
			p: s2.LatLng{Lat: math.Pi / 2, Lng: 0},
			n: s2.LatLng{Lat: -math.Pi / 6, Lng: 2 * math.Pi / 3},
			gm: &GeneralizedMercator{
				pos: r3.Vector{0, 0, 1},
				neg: r3.Vector{-math.Sqrt(3) / 4, 0.75, -0.5},
				i:   r3.Vector{-math.Sqrt(3) / 4, 0.75, 0.5},
				j:   r3.Vector{-math.Sqrt(3) / 2, -0.5, 0},
				k:   r3.Vector{0.25, -math.Sqrt(3) / 4, math.Sqrt(3) / 2},
				t:   2,
			},
		},
		{
			p: s2.LatLng{Lat: -math.Pi / 4, Lng: 0},
			n: s2.LatLng{Lat: -math.Pi / 4, Lng: math.Pi},
			gm: &GeneralizedMercator{
				pos: r3.Vector{math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				neg: r3.Vector{-math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
				i:   r3.Vector{0, 0, -1},
				j:   r3.Vector{0, 1, 0},
				k:   r3.Vector{1, 0, 0},
				t:   math.Sqrt2,
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
			t:   math.Inf(1),
		},
		ps: []proj{
			{s2.LatLng{Lat: math.Pi / 2}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, math.Inf(1)).X, Y: p0(0.5, math.Inf(1)).Y, Z: p0(0.5, math.Inf(1)).Z}}), r2.Point{X: 0, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, math.Inf(1)).X, Y: p1(0.5, math.Inf(1)).Y, Z: p1(0.5, math.Inf(1)).Z}}), r2.Point{X: math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, math.Inf(1)).X, Y: p2(0.5, math.Inf(1)).Y, Z: p2(0.5, math.Inf(1)).Z}}), r2.Point{X: math.Pi, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, math.Inf(1)).X, Y: p3(0.5, math.Inf(1)).Y, Z: p3(0.5, math.Inf(1)).Z}}), r2.Point{X: -math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLng{Lat: 0, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: math.Pi / 2}, r2.Point{X: math.Pi / 2, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: math.Pi}, r2.Point{X: math.Pi, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: -math.Pi / 2}, r2.Point{X: -math.Pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, math.Inf(1)).X, Y: p0(-0.5, math.Inf(1)).Y, Z: p0(-0.5, math.Inf(1)).Z}}), r2.Point{X: 0, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, math.Inf(1)).X, Y: p1(-0.5, math.Inf(1)).Y, Z: p1(-0.5, math.Inf(1)).Z}}), r2.Point{X: math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, math.Inf(1)).X, Y: p2(-0.5, math.Inf(1)).Y, Z: p2(-0.5, math.Inf(1)).Z}}), r2.Point{X: math.Pi, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, math.Inf(1)).X, Y: p3(-0.5, math.Inf(1)).Y, Z: p3(-0.5, math.Inf(1)).Z}}), r2.Point{X: -math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLng{Lat: -math.Pi / 2}, r2.Point{Y: math.Inf(-1)}},
		},
	},
	{
		gm: &GeneralizedMercator{
			pos: r3.Vector{0.5, -math.Sqrt(3) / 2, 0},
			neg: r3.Vector{0.5, math.Sqrt(3) / 2, 0},
			i:   r3.Vector{1, 0, 0},
			j:   r3.Vector{0, 0, 1},
			k:   r3.Vector{0, -1, 0},
			t:   2,
		},
		ps: []proj{
			{s2.LatLng{Lat: 0, Lng: -math.Pi / 3}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, 2).X, Y: -p0(0.5, 2).Z, Z: p0(0.5, 2).Y}}), r2.Point{X: 0, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, 2).X, Y: -p1(0.5, 2).Z, Z: p1(0.5, 2).Y}}), r2.Point{X: math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, 2).X, Y: -p2(0.5, 2).Z, Z: p2(0.5, 2).Y}}), r2.Point{X: math.Pi, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, 2).X, Y: -p3(0.5, 2).Z, Z: p3(0.5, 2).Y}}), r2.Point{X: -math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLng{Lat: 0, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: math.Pi / 2}, r2.Point{X: math.Pi / 2, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: math.Pi}, r2.Point{X: math.Pi, Y: 0}},
			{s2.LatLng{Lat: -math.Pi / 2}, r2.Point{X: -math.Pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, 2).X, Y: -p0(-0.5, 2).Z, Z: p0(-0.5, 2).Y}}), r2.Point{X: 0, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, 2).X, Y: -p1(-0.5, 2).Z, Z: p1(-0.5, 2).Y}}), r2.Point{X: math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, 2).X, Y: -p2(-0.5, 2).Z, Z: p2(-0.5, 2).Y}}), r2.Point{X: math.Pi, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, 2).X, Y: -p3(-0.5, 2).Z, Z: p3(-0.5, 2).Y}}), r2.Point{X: -math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLng{Lat: 0, Lng: math.Pi / 3}, r2.Point{Y: math.Inf(-1)}},
		},
	},
	{
		gm: &GeneralizedMercator{
			pos: r3.Vector{math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
			neg: r3.Vector{-math.Sqrt2 / 2, 0, -math.Sqrt2 / 2},
			i:   r3.Vector{0, 0, -1},
			j:   r3.Vector{0, 1, 0},
			k:   r3.Vector{1, 0, 0},
			t:   math.Sqrt2,
		},
		ps: []proj{
			{s2.LatLng{Lat: -math.Pi / 4, Lng: 0}, r2.Point{Y: math.Inf(1)}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(0.5, math.Sqrt2).Z, Y: p0(0.5, math.Sqrt2).Y, Z: -p0(0.5, math.Sqrt2).X}}), r2.Point{X: 0, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(0.5, math.Sqrt2).Z, Y: p1(0.5, math.Sqrt2).Y, Z: -p1(0.5, math.Sqrt2).X}}), r2.Point{X: math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(0.5, math.Sqrt2).Z, Y: p2(0.5, math.Sqrt2).Y, Z: -p2(0.5, math.Sqrt2).X}}), r2.Point{X: math.Pi, Y: math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(0.5, math.Sqrt2).Z, Y: p3(0.5, math.Sqrt2).Y, Z: -p3(0.5, math.Sqrt2).X}}), r2.Point{X: -math.Pi / 2, Y: math.Log(3) / 2}},
			{s2.LatLng{Lat: -math.Pi / 2, Lng: 0}, r2.Point{X: 0, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: math.Pi / 2}, r2.Point{X: math.Pi / 2, Y: 0}},
			{s2.LatLng{Lat: math.Pi / 2, Lng: 0}, r2.Point{X: math.Pi, Y: 0}},
			{s2.LatLng{Lat: 0, Lng: -math.Pi / 2}, r2.Point{X: -math.Pi / 2, Y: 0}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p0(-0.5, math.Sqrt2).Z, Y: p0(-0.5, math.Sqrt2).Y, Z: -p0(-0.5, math.Sqrt2).X}}), r2.Point{X: 0, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p1(-0.5, math.Sqrt2).Z, Y: p1(-0.5, math.Sqrt2).Y, Z: -p1(-0.5, math.Sqrt2).X}}), r2.Point{X: math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p2(-0.5, math.Sqrt2).Z, Y: p2(-0.5, math.Sqrt2).Y, Z: -p2(-0.5, math.Sqrt2).X}}), r2.Point{X: math.Pi, Y: -math.Log(3) / 2}},
			{s2.LatLngFromPoint(s2.Point{r3.Vector{X: p3(-0.5, math.Sqrt2).Z, Y: p3(-0.5, math.Sqrt2).Y, Z: -p3(-0.5, math.Sqrt2).X}}), r2.Point{X: -math.Pi / 2, Y: -math.Log(3) / 2}},
			{s2.LatLng{Lat: -math.Pi / 4, Lng: math.Pi}, r2.Point{Y: math.Inf(-1)}},
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
		(a.t == b.t || math.Abs(a.t-b.t) < 1e-15)
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
