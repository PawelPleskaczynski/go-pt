package main

import (
	"math"
)

// based on https://www.scratchapixel.com/lessons/procedural-generation-virtual-worlds/simulating-sky/simulating-colors-of-the-sky

type Atmosphere struct {
	sunDirection        Tuple
	sunSize             float64
	sunColor            Color
	earthRadius         float64
	atmosphereRadius    float64
	rayleighScaleHeight float64
	mieScaleHeight      float64
	betaR               Tuple
	betaM               Tuple
	mulR                float64
	mulM                float64
	mulAll              float64
}

func NewEarthAtmosphere(sunDirection Tuple) Atmosphere {
	atm := Atmosphere{}
	atm.sunDirection = sunDirection.Normalize()
	atm.sunSize = 2 * math.Pi / 180 // 2 degrees is bigger than in real life, but for now there's no importance sampling so it converges really slow on 0.5 deg Sun
	atm.sunColor = Color{200, 200, 200}
	atm.earthRadius = 6360e3
	atm.atmosphereRadius = 6420e3
	atm.rayleighScaleHeight = 7994
	atm.mieScaleHeight = 1200
	atm.betaR = Tuple{3.8e-6, 13.5e-6, 33.1e-6, 0}
	atm.betaM = Tuple{21e-6, 21e-6, 21e-6, 0}
	atm.mulR = 1.0
	atm.mulM = 10.0
	atm.mulAll = 20.0
	return atm
}

func skyHit(origin, direction Tuple, radius, tMin, tMax float64) (bool, float64, float64) {
	a := direction.Dot(direction)
	b := 2.0 * origin.Dot(direction)
	c := origin.Dot(origin) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant > 0.0 {
		t0 := (-b - math.Sqrt(discriminant)) / (2.0 * a)
		t1 := (-b + math.Sqrt(discriminant)) / (2.0 * a)

		if t0 > t1 {
			return true, t1, t0
		}
		return true, t0, t1
	}
	return false, -math.MaxFloat64, -math.MaxFloat64
}

func (atm Atmosphere) ComputeIncidentLight(origin, direction Tuple, tMin, tMax float64) Tuple {
	hitResult, t0, t1 := skyHit(origin, direction, atm.atmosphereRadius, tMin, tMax)

	if !hitResult || t1 < 0 {
		return Tuple{0.1, 0, 0, 0}
	}

	if t0 > tMin && t0 > 0.0 {
		tMin = t0
	}

	if t1 < tMax {
		tMax = t1
	}

	// check if the sun was hit
	sunHit := false
	sunCenter := atm.sunDirection
	distance := math.Sqrt(math.Pow(direction.x-sunCenter.x, 2) + math.Pow(direction.y-sunCenter.y, 2) + math.Pow(direction.z-sunCenter.z, 2))

	if distance <= atm.sunSize {
		sunHit = true
	}

	numSamples := 16
	numSamplesLight := 8
	segmentLength := (tMax - tMin) / float64(numSamples)

	tCurrent := tMin
	sumR, sumM := Tuple{0, 0, 0, 0}, Tuple{0, 0, 0, 0}
	opticalDepthR, opticalDepthM := 0.0, 0.0

	mu := direction.Dot(atm.sunDirection)
	phaseR := 3.0 / (16.0 * math.Pi) * (1.0 + mu*mu)
	g := 0.76
	phaseM := 3.0 / (8.0 * math.Pi) * ((1.0 - g*g) * (1.0 + mu*mu)) / ((2.0 + g*g) * math.Pow(1.0+g*g-2.0*g*mu, 1.5))

	for i := 0; i < numSamples; i++ {
		samplePosition := origin.Add(direction.MulScalar(tCurrent + segmentLength*0.5))
		height := samplePosition.Magnitude() - atm.earthRadius

		hr := math.Exp(-height/atm.rayleighScaleHeight) * segmentLength
		hm := math.Exp(-height/atm.mieScaleHeight) * segmentLength

		opticalDepthR += hr
		opticalDepthM += hm

		_, _, t1Light := skyHit(samplePosition, atm.sunDirection, atm.atmosphereRadius, tMin, tMax)
		segmentLengthLight := t1Light / float64(numSamplesLight)
		tCurrentLight, opticalDepthLightR, opticalDepthLightM := 0.0, 0.0, 0.0
		var j int
		for j = 0; j < numSamplesLight; j++ {
			samplePositionLight := samplePosition.Add((atm.sunDirection).AddScalar(tCurrentLight + segmentLengthLight*0.5))
			heightLight := samplePositionLight.Magnitude() - atm.earthRadius
			if heightLight < 0 {
				break
			}
			opticalDepthLightR += math.Exp(-heightLight/atm.rayleighScaleHeight) * segmentLengthLight
			opticalDepthLightM += math.Exp(-heightLight/atm.mieScaleHeight) * segmentLengthLight
			tCurrentLight += segmentLengthLight
		}

		if j == numSamplesLight {
			tau := (atm.betaR.MulScalar(opticalDepthR + opticalDepthLightR)).Add(atm.betaM.MulScalar(1.1).MulScalar(opticalDepthM + opticalDepthLightM))
			attenuation := Tuple{math.Exp(-tau.x), math.Exp(-tau.y), math.Exp(-tau.z), 0}
			sumR = sumR.Add(attenuation.MulScalar(hr))
			sumM = sumM.Add(attenuation.MulScalar(hm))
		}
		tCurrent += segmentLength
	}

	retColor := (((sumR.Mul(atm.betaR)).MulScalar(phaseR).MulScalar(atm.mulR)).Add((sumM.Mul(atm.betaM)).MulScalar(phaseM).MulScalar(atm.mulM))).MulScalar(atm.mulAll)

	if sunHit {
		retColor = retColor.Mul(Tuple{atm.sunColor.r, atm.sunColor.g, atm.sunColor.b, 0})
	}

	return retColor
}
