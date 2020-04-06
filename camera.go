package main

import (
	"math"
	"math/rand"
)

type Camera struct {
	origin, lowerLeftCorner, horizontal, vertical, u, v, w Tuple
	lensRadius                                             float64
}

func randomDisk(generator rand.Rand) Tuple {
	var p Tuple
	for {
		p = Tuple{RandFloat(generator), RandFloat(generator), 0, 0}.MulScalar(2.0).Subtract(Tuple{1, 1, 0, 0})
		if p.Dot(p) <= 1.0 {
			break
		}
	}
	return p
}

func getCamera(lookFrom, lookAt, up Tuple, fLength, aspect, fNumber, focusDistance float64) Camera {
	var c Camera
	aperture := (fLength / 1000) / fNumber
	c.lensRadius = aperture / 2
	theta := 2 * math.Atan(24/(2*fLength))
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight
	c.origin = lookFrom
	c.w = (lookFrom.Subtract(lookAt)).Normalize()
	c.u = up.Cross(c.w).Normalize()
	c.v = c.w.Cross(c.u)
	c.lowerLeftCorner = c.origin.Add(c.u.MulScalar(-(halfWidth) * focusDistance)).Add(c.v.MulScalar(-(halfHeight) * focusDistance)).Add(c.w.Negate().MulScalar(focusDistance))
	c.horizontal = c.u.MulScalar(2 * halfWidth * focusDistance)
	c.vertical = c.v.MulScalar(2 * halfHeight * focusDistance)

	return c
}

func (c Camera) getRay(s, t float64, generator rand.Rand) Ray {
	randomDisk := randomDisk(generator).MulScalar(c.lensRadius)
	offset := c.u.MulScalar(randomDisk.x).Add(c.v.MulScalar(randomDisk.y))
	return Ray{c.origin.Add(offset), c.lowerLeftCorner.Add(c.horizontal.MulScalar(s)).Add(c.vertical.MulScalar(t)).Subtract(c.origin).Subtract(offset)}
}
