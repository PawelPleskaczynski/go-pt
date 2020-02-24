package main

import (
	"math/rand"
)

func RandFloat(generator rand.Rand) float64 {
	return generator.Float64()
}

func RandInUnitSphere(generator rand.Rand) Tuple {
	p := Tuple{0, 0, 0, 0}
	for {
		p = (Tuple{RandFloat(generator), RandFloat(generator), RandFloat(generator), 1}.Subtract(Tuple{0.5, 0.5, 0.5, 0})).MulScalar(2.0)
		if p.Magnitude()*p.Magnitude() < 1.0 {
			break
		}
	}
	return p.Normalize()
}

func RandInUnitHemisphere(generator rand.Rand, normal Tuple) Tuple {
	p := Tuple{0, 0, 0, 0}
	for {
		p = (Tuple{RandFloat(generator), RandFloat(generator), RandFloat(generator), 1}.Subtract(Tuple{0.5, 0.5, 0.5, 0})).MulScalar(2.0)
		if p.Dot(normal)/(p.Magnitude()*normal.Magnitude()) >= 0.0 {
			break
		}
	}
	return p
}
