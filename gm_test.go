package gm

import (
	"math"
	"testing"

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

func gmApproxEqual(a, b *GeneralizedMercator) bool {
	return approxEqual(a.pos, b.pos) &&
		approxEqual(a.neg, b.neg) &&
		approxEqual(a.i, b.i) &&
		approxEqual(a.j, b.j) &&
		approxEqual(a.k, b.k) &&
		(a.t == b.t || math.Abs(a.t-b.t) < 1e-15)
}
