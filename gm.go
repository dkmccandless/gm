/*
Package gm implements the generalized spherical Mercator projection.

The generalized Mercator projection maps a spherical surface onto a flat plane with respect to two poles,
which must be distinct but need not be antipodes.

Like the Mercator projection, it is finite in width but infinite in height: the x coordinate of the projection
corresponds to an analogue of longitude along the great circle of points equidistant from the poles,
and the y coordinate is a function of distance from this great circle toward either pole (in analogy with latitude);
the poles themselves have infinite y values.

The Mercator projection is the special case corresponding to poles at latitude ±90°.
*/
package gm

import (
	"math"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// GeneralizedMercator defines the generalized spherical Mercator projection with poles at pos and neg.
// Initialize a new GeneralizedMercator with New. The zero value of type GeneralizedMercator will result in undefined behavior.
type GeneralizedMercator struct {
	// To explicitly reflect that Pos and Neg are vector quantities, their names are capitalized in the documentation.
	// However, the field identifiers themselves are not capitalized so that they remain unexported.

	// Pos and Neg are unit vectors representing the poles of the projection.
	pos, neg r3.Vector

	// i, j, and k are a convenient orthonormal basis for the projection operations.
	i, j, k r3.Vector

	// The vector T, discussed below, denotes the closest point to the origin on
	// the line of intersection of the planes tangent to the unit sphere at Pos and Neg.
	// t is the magnitude of T.

	// t is the (possibly infinite) distance to the line of intersection
	// of the planes tangent to the unit sphere at Pos and Neg.
	t float64
}

// New returns a pointer to a GeneralizedMercator with poles at pos and neg.
// It panics if pos and neg are equal.
func New(pos, neg s2.LatLng) *GeneralizedMercator {
	gm := &GeneralizedMercator{
		pos: s2.PointFromLatLng(pos).Vector,
		neg: s2.PointFromLatLng(neg).Vector,
	}

	// Snap to nearest integer if necessary to avoid math.Cos rounding error
	for _, v := range []r3.Vector{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {-1, 0, 0}, {0, -1, 0}, {0, 0, -1}} {
		if approxEqual(gm.pos, v) {
			gm.pos = v
		}
		if approxEqual(gm.neg, v) {
			gm.neg = v
		}
	}
	if approxEqual(gm.pos, gm.neg) {
		panic("indistinguishable poles")
	}

	// The right-handed orthonormal basis (i, j, k) is defined such that
	// k is in the direction from Neg to Pos, and i is on the line bisecting Pos and Neg.
	// If Pos and Neg are not antipodes, i points in the direction of their average,
	// or else in the direction of the intersection between their bisector and the prime meridian
	// (and the Equator at 0° longitude, if their bisector contains the prime meridian,
	// or the north pole, if their bisector contains both poles), with T on the positive i axis.
	//
	// The projection is defined according to this basis to be centered at
	// the intersection of the positive i axis with the reference sphere,
	// and oriented such that, at that point, j points to the right and k points up;
	// Pos lies at the (infinitely distant) top of the projection, and Neg at the bottom.
	gm.k = gm.pos.Sub(gm.neg).Normalize()

	switch {
	case approxEqual(gm.pos, gm.neg.Mul(-1)):
		// Pos and Neg are antipodes; their tangent planes intersect at the line at infinity.
		gm.t = math.Inf(1)

		// Define the basis such that i lies on the great circle equidistant from them and on the prime meridian:
		// at 0° latitude, if the great circle contains it, or else at the north pole,
		// if the great circle contains it, or else at their unique intersection point.
		switch {
		case gm.pos.X == 0:
			// If Pos and Neg are on the great circle of longitude ±90°,
			// i lies at the intersection of the prime meridian with the Equator.
			gm.i = r3.Vector{1, 0, 0}

		case gm.pos.Z == 0:
			// If Pos and Neg are elsewhere on the equator, i lies at the north pole.
			gm.i = r3.Vector{0, 0, 1}

		default:
			// The great circle equidistant from Pos and Neg intersects the prime meridian at a single point; i lies at that point.
			// Its latitude is φ_i = arctan(-cos(λ_P)/tan(φ_P)), the solution of P • i == 0, λ_i == 0,
			// where P is a pole with non-negative latitude so that φ is in the range [-π/2, π/2].
			p := pos
			if p.Lat < 0 {
				p = neg
			}
			phi := math.Atan2(-math.Cos(float64(p.Lng)), math.Tan(float64(p.Lat)))
			gm.i = r3.Vector{math.Cos(phi), 0, math.Sin(phi)}
		}

	default:
		// Pos and Neg are not antipodes; define the basis such that i lies at the closest point equidistant from them.
		gm.i = gm.pos.Add(gm.neg).Normalize()

		// T = (Pos + Neg) / (1 + Pos•Neg) is the solution of T • Pos == 1, T • Neg == 1, and T • (Pos × Neg) == 0.
		gm.t = gm.pos.Add(gm.neg).Norm() / (1 + gm.pos.Dot(gm.neg))
	}

	// j is orthogonal to Pos and Neg in the direction of increasing projectional longitude at the zero point.
	// If Pos and Neg are not antipodes, the intersection line of the planes tangent to the unit sphere at Pos and Neg is parallel to j.
	gm.j = gm.k.Cross(gm.i)

	return gm
}

// Project converts ll to a projected 2D point.
func (gm *GeneralizedMercator) Project(ll s2.LatLng) r2.Point {
	P := s2.PointFromLatLng(ll).Vector
	switch {
	case approxEqual(P, gm.pos):
		return r2.Point{Y: math.Inf(1)}
	case approxEqual(P, gm.neg):
		return r2.Point{Y: math.Inf(-1)}
	}

	var (
		beta   = math.Copysign(float64(gm.i.Sub(P.Mul(1/gm.t)).Cross(gm.j).Angle(gm.k)), P.Dot(gm.k))
		iprime = s2.Rotate(s2.Point{gm.i}, s2.Point{gm.j}, s1.Angle(beta)).Vector
		kprime = s2.Rotate(s2.Point{gm.k}, s2.Point{gm.j}, s1.Angle(beta)).Vector

		C = kprime.Mul(P.Dot(kprime))

		// The vertical component of the projection is the same function of angle OPC as the Mercator projection,
		// the difference being that this angle is in general no longer equal to the point's latitude.
		// To differentiate the generalized case, call this angle ψ.
		psi = math.Copysign(float64(P.Angle(P.Sub(C))), P.Dot(gm.k))
		z   = math.Log(math.Tan(math.Pi/4 + psi/2))

		// The horizontal component is equal to the measure of angle TCP.
		theta = math.Copysign(float64(iprime.Angle(P.Sub(C))), P.Sub(C).Dot(gm.j))
	)

	return r2.Point{X: theta, Y: z}
}

// Unproject converts a projected point p to a location on the reference sphere.
func (gm *GeneralizedMercator) Unproject(p r2.Point) s2.LatLng {
	switch {
	case math.IsInf(p.Y, 1):
		return s2.LatLngFromPoint(s2.Point{gm.pos})
	case math.IsInf(p.Y, -1):
		return s2.LatLngFromPoint(s2.Point{gm.neg})
	}

	// The locus of points on the sphere with equal vertical projective coordinate is a circle; let C be its center.
	// The vertical component of the projection is the same function of angle OPC as the Mercator projection,
	// the difference being that this angle is in general no longer equal to the point's latitude.
	// To differentiate the generalized case, call this angle ψ.
	// Use the inverse function to construct C.
	var (
		psi = 2*math.Atan(math.Exp(p.Y)) - math.Pi/2

		// β is the measure of angle OTC, related by sin(β) = sin(ψ)/|T|.
		// This is also the angular distance of C from the positive k axis.
		beta = math.Asin(math.Sin(psi) / gm.t)

		// Rotate the basis by an angle of β around the j axis such that C is parallel to the k' axis,
		// and T-C is parallel to the i' axis.
		iprime = s2.Rotate(s2.Point{gm.i}, s2.Point{gm.j}, s1.Angle(beta))
		kprime = s2.Rotate(s2.Point{gm.k}, s2.Point{gm.j}, s1.Angle(beta))

		C = kprime.Mul(math.Sin(psi))

		// P is the point on the circle centered at C such that theta is the measure of angle TCP.
		P = s2.Rotate(iprime, kprime, s1.Angle(p.X)).Mul(math.Cos(psi)).Add(C)
	)

	return s2.LatLngFromPoint(s2.Point{P})
}

// approxEqual is equivalent to r3.Vector's ApproxEqual method but with a larger tolerance.
func approxEqual(a, b r3.Vector) bool {
	// r3's epsilon of 1e-16 is too strict to accommodate some values returned by s2.PointFromLatLng
	// due to propagation of the float64 rounding error math.Cos(math.Pi/2) == 6.123233995736757e-17.
	// For example, s2.PointFromLatLng(Lat: 0, Lng: math.Pi) == (-1, -1.2246467991473515e-16, 0),
	// and s2.LatLng{Lat: math.Pi/2}.ApproxEqual(s2.LatLng{Lat: math.Pi/2, Lng: math.Pi}) is false.
	// 1e-15 is still only about 6.4 nanometers at the Earth's surface.
	const epsilon = 1e-15
	return math.Abs(a.X-b.X) < epsilon && math.Abs(a.Y-b.Y) < epsilon && math.Abs(a.Z-b.Z) < epsilon
}
