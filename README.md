# gm

Package gm implements the generalized spherical Mercator projection.

## Install

	go get github.com/dkmccandless/gm

## Overview

The generalized Mercator projection maps a spherical surface onto a flat plane with respect to two poles, which must be distinct but need not be antipodes.

Like the [Mercator projection](https://en.wikipedia.org/wiki/Mercator_projection), it is finite in width but infinite in height: the x coordinate of the projection corresponds to an analogue of longitude along the great circle of points equidistant from the poles, and the y coordinate is a function of distance from this great circle toward either pole (in analogy with latitude); the poles themselves have infinite y values.

The Mercator projection is the special case corresponding to poles at latitude ±90°.

## Definition of the projection

Given a "positive" and "negative" pole described by unit vectors **P<sub>+</sub>** and **P<sub>-</sub>**, let *α* be the half-angle between them. Then the distance from the origin to the line of intersection of the planes tangent to the unit circle at **P<sub>+</sub>** and **P<sub>-</sub>** is *d* = sec(*α*). The unit sphere can be parametrized according to the sheaf of planes passing through this line. (In the case that the poles are antipodes, this line is the line at infinity, and these planes are all parallel; in the case that the poles of the projection are the north and south pole, these are the parallels of latitude.)

Define the right-handed orthonormal basis (**î**, **ĵ**, **k̂**) as follows. Let **k̂** be a unit vector in the direction of **P<sub>+</sub>** - **P<sub>-</sub>**. If **P<sub>+</sub>** and **P<sub>-</sub>** are not antipodes, then let **ĵ** be a unit vector in the direction of **P<sub>+</sub>** × **P<sub>-</sub>**, which is parallel to the line of intersection between the planes tangent to the unit circle at **P<sub>+</sub>** and **P<sub>-</sub>**. Then **î** points toward the closest point on this line to the origin. In terms of the poles, **î** = (**P<sub>+</sub>** + **P<sub>-</sub>**) / 2·cos(*α*), **ĵ** = **P<sub>+</sub>** × **P<sub>-</sub>** / 2·sin(*α*)cos(*α*), and **k̂** = (**P<sub>+</sub>** - **P<sub>-</sub>**) / 2·sin(*α*).

If **P<sub>+</sub>** and **P<sub>-</sub>** are antipodes, on the other hand, the tangent planes are parallel and their line of intersection is the line at infinity, and so it is necessary to choose the direction of the **î** axis. Call the line equidistant from the poles the *projective equator* (not capitalized). Choose **î** to lie on the intersection of the projective equator with the prime meridian (with **ĵ** subsequently determined by the right-hand rule), according to the following cases:

* If the poles are on the (terrestrial) Equator at longitude ±90°, the projective equator contains the entire prime meridian; define **î** to be the point on the Equator at 0° longitude.
* If the poles are antipodes elsewhere on the Equator, the projective equator contains the north and south poles; define **î** to pass through the north pole.
* If the poles are antipodes not on the Equator, the projective equator intersects the prime meridian at a single point; define **î** to pass through that point.

Given a point **P** on the unit sphere at latitude-longitude coordinates (*φ*, *θ*), there is a unique plane containing both **P** and the line of intersection of the tangent planes. (**î** - **P**/*d*) × **ĵ** is a normal vector to this plane; let *β* be its angle with respect to **k̂** (equivalently, the dihedral angle of this plane with respect to the plane of the projective equator). It is convenient to consider a modified basis rotated around **ĵ** by this angle, such that **î′** = **î**·cos(*β*) - **k̂**·sin(*β*), **ĵ′** = **ĵ**, and **k̂′** = **î**·sin(*β*) + **k̂**·cos(*β*). This plane's intersection with the unit sphere describes a circle parallel to the **î′ĵ′**-plane. Let *ψ* be the latitude of the circle with respect to the **î′ĵ′**-plane, related by sin(*ψ*)/*d* = sin(*β*). The value of *ψ* ranges from -*π*/2 at the negative pole, through 0 on the projective equator, to *π*/2 at the positive pole. *ψ* is the generalized analogue of the conventional latitude coordinate *φ*, and the vertical projective coordinate *y* is related by the same function of *ψ* as in the Mercator projection: *y* = ln(tan(*π*/4 + *ψ*/2)), and inversely *ψ* = 2·arctan(*e*<sup>*y*</sup>) - *π*/2.

The horizontal projective coordinate is defined as the angular distance of the projection of **P** in the **î′ĵ′**-plane from the **î′** axis—the longitude of **P** in the **î′ĵ′k̂′** basis: *x* = atan2(**P** · **ĵ′**, **P** · **î′**).
