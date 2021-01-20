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

/*
The projection operations are expressed in terms of the right-handed orthonormal basis (i, j, k), defined as follows.
The k axis is parallel to the vector from Neg to Pos, and the i and j axes are on the great circle bisecting them.
If Pos and Neg are not antipodes, then the i axis passes through the midpoint of the shorter leg of the great circle
containing them. Otherwise, the i axis is defined according to the intersection of the great circle bisecting them with
the prime meridian:

  1. If the bisector contains the prime meridian, the i axis intersects the prime Meridian at the Equator.
  2. If the bisector contains the North and South Poles, the i axis passes through the North Pole.
  3. If the bisector intersects the prime meridian at a single point, the i axis passes through that point.

Equivalently, the i axis passes through the intersection of the prime meridian with the Equator, if this point is
equidistant from Pos and Neg, or else the North Pole, if it is equidistant from Pos and Neg, or else the unique point
on the prime meridian that is equidistant from Pos and Neg.

In this basis, the intersection line of the planes tangent to the unit sphere at Pos and Neg is described by the set of
points with i == t, k == 0. For any point P on the unit sphere, consider the unique plane containing both this line and
P. Let β be the dihedral angle between this plane and the ij-plane. The intersection of this plane with the unit sphere
is a circle; let C be its center and ψ be the complement of its polar angle, related by |C| = sin(ψ) = t*sin(β). ψ is
the generalized analogue of the conventional latitude coordinate φ, and the vertical projective coordinate y is related
by the same function of ψ as in the Mercator projection: y = ln(tan(π/4 + ψ/2)), and inversely ψ = 2*arctan(e^y) - π/2.
The horizontal projective coordinate x is equal to the arc length along the circle from the point on the circle with
maximal i coordinate to P.

It is convenient to define a new basis in which to consider P, rotated by an angle β around the j axis such that C lies
on the k' axis. In this basis, ψ is simply the latitude of P, and x is P's longitude from the positive i' axis.
*/

// New returns a pointer to a GeneralizedMercator with poles at pos and neg.
// It panics if pos and neg are equal.
func New(pos, neg s2.LatLng) *GeneralizedMercator {
	gm := &GeneralizedMercator{
		// Snap each coordinate to the nearest integer if necessary to avoid math.Cos rounding error
		pos: snapToInts(s2.PointFromLatLng(pos).Vector),
		neg: snapToInts(s2.PointFromLatLng(neg).Vector),
	}

	if approxEqual(gm.pos, gm.neg) {
		panic("indistinguishable poles")
	}

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
			// the i axis passes through the intersection of the prime meridian with the Equator.
			gm.i = r3.Vector{1, 0, 0}

		case gm.pos.Z == 0:
			// If Pos and Neg are elsewhere on the equator, the i axis passes through the north pole.
			gm.i = r3.Vector{0, 0, 1}

		default:
			// The great circle equidistant from Pos and Neg intersects the prime meridian at a single point;
			// the i axis intersects that point.
			gm.i = r3.Vector{0, gm.pos.Z, 0}.Cross(gm.pos).Normalize()
		}

	default:
		// Pos and Neg are not antipodes; the i axis passes through the closest point equidistant from them.
		gm.i = gm.pos.Add(gm.neg).Normalize()

		// T = (Pos + Neg) / (1 + Pos•Neg) is the solution of T • Pos == 1, T • Neg == 1, and T • (Pos × Neg) == 0.
		// t is the magnitude of T, equivalent to the secant of the half-angle between Pos and Neg.
		gm.t = gm.pos.Add(gm.neg).Norm() / (1 + gm.pos.Dot(gm.neg))
	}

	// j is orthogonal to Pos and Neg in the direction of increasing projectional longitude at the zero point.
	// If Pos and Neg are not antipodes, the intersection line of the planes tangent to the unit sphere at Pos and Neg is parallel to the j axis.
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

		psi = math.Copysign(float64(P.Angle(P.Sub(C))), P.Dot(gm.k))
		y   = math.Log(math.Tan(math.Pi/4 + psi/2))
		x   = math.Copysign(float64(iprime.Angle(P.Sub(C))), P.Sub(C).Dot(gm.j))
	)

	return r2.Point{x, y}
}

// Unproject converts a projected point p to a location on the reference sphere.
func (gm *GeneralizedMercator) Unproject(p r2.Point) s2.LatLng {
	switch {
	case math.IsInf(p.Y, 1):
		return s2.LatLngFromPoint(s2.Point{gm.pos})
	case math.IsInf(p.Y, -1):
		return s2.LatLngFromPoint(s2.Point{gm.neg})
	}

	var (
		psi  = 2*math.Atan(math.Exp(p.Y)) - math.Pi/2
		beta = math.Asin(math.Sin(psi) / gm.t)

		iprime = s2.Rotate(s2.Point{gm.i}, s2.Point{gm.j}, s1.Angle(beta))
		kprime = s2.Rotate(s2.Point{gm.k}, s2.Point{gm.j}, s1.Angle(beta))

		C = kprime.Mul(math.Sin(psi))

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

// snapToInts returns v with any component approximately equal to an integer rounded to that integer.
func snapToInts(v r3.Vector) r3.Vector {
	const epsilon = 1e-15
	if r := math.Round(v.X); math.Abs(v.X-r) < epsilon {
		v.X = r
	}
	if r := math.Round(v.Y); math.Abs(v.Y-r) < epsilon {
		v.Y = r
	}
	if r := math.Round(v.Z); math.Abs(v.Z-r) < epsilon {
		v.Z = r
	}
	return v
}
