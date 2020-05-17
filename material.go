package main

import (
	"math"
	"math/rand"
)

const (
	BSDF = iota
	Lambertian
	Emission
)

type Material struct {
	material           int
	albedo             Texture
	roughness          float64
	ior                float64
	clearcoat          float64
	clearcoatRoughness float64
	metalicity         float64
	transmission       float64
}

func getLambertian(albedo Texture) Material {
	return Material{Lambertian, albedo, 0, 1.5, 0, 0, 0, 0}
}

func getGlossy(albedo Texture, roughness, clearcoat float64) Material {
	return Material{BSDF, albedo, roughness, 1.5, clearcoat, roughness, 0, 0}
}

func getDielectric(albedo Texture, roughness, clearcoat, ior float64) Material {
	return Material{BSDF, albedo, roughness, ior, clearcoat, roughness, 0, 1}
}

func getMetal(albedo Texture, roughness, clearcoat, clearcoatRoughness float64) Material {
	return Material{BSDF, albedo, roughness, 1.5, clearcoat, clearcoatRoughness, 1, 0}
}

func getEmission(albedo Texture) Material {
	return Material{Emission, albedo, 0, 0, 0, 0, 0, 0}
}

func sampleGGX(xi1, xi2, a float64) (float64, float64) {
	phi := 2.0 * math.Pi * xi1
	theta := math.Acos(math.Sqrt((1.0 - xi2) / ((a*a-1.0)*xi2 + 1.0)))
	return phi, theta
}

func generateGGXNormal(normal Tuple, roughness float64, generator rand.Rand) (Tuple, float64) {
	theta, phi := sampleGGX(RandFloat(generator), RandFloat(generator), roughness*roughness)

	x := phi * math.Cos(theta)
	y := phi * math.Sin(theta)
	z := math.Sqrt(1 - x*x - y*y)

	ggxNormal := Tuple{x, y, z, 0}
	uvw := buildFromW(normal)
	return uvw.local(ggxNormal).Normalize(), phi
}

func (m Material) Scatter(r Ray, rec HitRecord, attenuation *Color, scattered *Ray, generator rand.Rand) bool {
	incoming := r.direction.Normalize()
	switch m.material {
	case Lambertian:
		target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
		*scattered = Ray{rec.p, target.Subtract(rec.p)}
		*attenuation = m.albedo.color(rec)
		return true
	case Emission:
		return true
	case BSDF:
		var outwardNormal Tuple
		var refracted Tuple

		var niOverNt float64
		var reflectProbability float64
		var cosine float64

		var n1, n2 float64
		var phi float64

		*attenuation = m.albedo.color(rec)

		clearcoatNormal, phi := generateGGXNormal(rec.normal, rec.material.clearcoatRoughness, generator)
		rec.normal, _ = generateGGXNormal(rec.normal, rec.material.roughness, generator)

		if incoming.Dot(rec.normal) > 0 {
			n1 = rec.material.ior
			n2 = 1.0
		} else {
			n1 = 1.0
			n2 = rec.material.ior
		}

		reflected := incoming.Reflection(rec.normal)
		reflectedClearcoat := incoming.Reflection(clearcoatNormal)
		clearcoat := math.Pow(rec.material.clearcoat, 1/2.2)
		reflectProbability = Fresnel(n1, n2, rec.normal, incoming) * (1 - phi)
		reflectProbability = clearcoat * reflectProbability

		if RandFloat(generator) <= m.metalicity {
			if RandFloat(generator) <= reflectProbability {
				*scattered = Ray{rec.p, reflectedClearcoat}
				*attenuation = Color{1, 1, 1}
			} else {
				*scattered = Ray{rec.p, reflected}
			}
			return (scattered.direction.Dot(rec.normal) > 0)
		}

		if RandFloat(generator) <= m.transmission {
			if incoming.Dot(rec.normal) > 0 {
				outwardNormal = rec.normal.MulScalar(-1)
				niOverNt = m.ior
				cosine = m.ior * incoming.Dot(rec.normal) / incoming.Magnitude()
			} else {
				outwardNormal = rec.normal
				niOverNt = 1.0 / m.ior
				cosine = -(incoming.Dot(rec.normal) / incoming.Magnitude())
			}

			if incoming.Refraction(outwardNormal, niOverNt, &refracted) {
				reflectProbability = Schlick(cosine, m.ior)
				reflectProbability = clearcoat * reflectProbability
				if reflectProbability > 1.0 {
					reflectProbability = 1.0
				}
			} else {
				reflectProbability = 1.0
			}

			if RandFloat(generator) <= reflectProbability {
				*scattered = Ray{rec.p, reflectedClearcoat}
				*attenuation = Color{1, 1, 1}
			} else {
				*scattered = Ray{rec.p, refracted}
			}

			return true
		}

		if RandFloat(generator) <= reflectProbability {
			*scattered = Ray{rec.p, reflectedClearcoat}
			*attenuation = Color{1, 1, 1}
		} else {
			target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
			*scattered = Ray{rec.p, target.Subtract(rec.p)}
		}

		return true
	default:
		return false
	}
}
