package main

import (
	"math"
	"math/rand"
)

const (
	BSDF = iota
	Lambertian
	Metal
	Dielectric
	Emission
)

type Material struct {
	material     int
	albedo       Texture
	roughness    float64
	ior          float64
	specularity  float64
	metalicity   float64
	transmission float64
}

func getLambertian(albedo Texture) Material {
	return Material{Lambertian, albedo, 0, 1.5, 0, 0, 0}
}

func getDiffuse(albedo Texture, roughness, specularity float64) Material {
	return Material{BSDF, albedo, roughness, 1.5, specularity, 0, 0}
}

func getDielectric(albedo Texture, roughness, specularity, ior float64) Material {
	return Material{BSDF, albedo, roughness, ior, specularity, 0, 1}
}

func getMetal(albedo Texture, roughness float64) Material {
	return Material{BSDF, albedo, roughness, 1.5, 0, 1, 0}
}

func getEmission(albedo Texture) Material {
	return Material{Emission, albedo, 0, 0, 0, 0, 0}
}

func sampleGGX(xi1, xi2, a float64) (float64, float64) {
	phi := 2.0 * math.Pi * xi1
	theta := math.Acos(math.Sqrt((1.0 - xi2) / ((a*a-1.0)*xi2 + 1.0)))
	return phi, theta
}

func (m Material) Scatter(r Ray, rec HitRecord, attenuation *Color, scattered *Ray, generator rand.Rand) bool {
	incoming := r.direction.Normalize()
	if m.material == Lambertian {
		target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
		*scattered = Ray{rec.p, target.Subtract(rec.p)}
		*attenuation = m.albedo.color(rec)
		return true
	} else if m.material == Metal {
		reflected := incoming.Reflection(rec.normal)
		*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
		*attenuation = m.albedo.color(rec)
		return (scattered.direction.Dot(rec.normal) > 0)
	} else if m.material == Dielectric {
		var outwardNormal Tuple
		var refracted Tuple

		var niOverNt float64
		var reflectProbability float64
		var cosine float64

		*attenuation = m.albedo.color(rec)
		reflected := incoming.Reflection(rec.normal)

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
			reflectProbability = Schlick(cosine, m.ior) + m.specularity/4
			if reflectProbability > 1.0 {
				reflectProbability = 1.0
			}
		} else {
			reflectProbability = 1.0
		}

		if RandFloat(generator) <= reflectProbability {
			*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
			*attenuation = Color{reflectProbability, reflectProbability, reflectProbability}
		} else {
			*scattered = Ray{rec.p, refracted.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
		}

		return true
	} else if m.material == Emission {
		return true
	} else if m.material == BSDF {
		var outwardNormal Tuple
		var refracted Tuple

		var niOverNt float64
		var reflectProbability float64
		var cosine float64

		*attenuation = m.albedo.color(rec)

		theta, phi := sampleGGX(RandFloat(generator), RandFloat(generator), m.roughness*m.roughness)

		x := phi * math.Cos(theta)
		y := phi * math.Sin(theta)

		z := math.Sqrt(1 - x*x - y*y)

		ggxNormal := Tuple{x, y, z, 0}
		uvw := buildFromW(rec.normal)
		rec.normal = uvw.local(ggxNormal).Normalize()

		reflected := incoming.Reflection(rec.normal)

		if RandFloat(generator) <= m.metalicity {
			*scattered = Ray{rec.p, reflected}
			return (scattered.direction.Dot(rec.normal) > 0)
		}

		specularity := rec.material.specularity / 4

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
				reflectProbability = (specularity + (1.0-specularity)*reflectProbability)
				if reflectProbability > 1.0 {
					reflectProbability = 1.0
				}
			} else {
				reflectProbability = 1.0
			}

			if RandFloat(generator) <= reflectProbability {
				*scattered = Ray{rec.p, reflected}
				*attenuation = Color{reflectProbability, reflectProbability, reflectProbability}
			} else {
				*scattered = Ray{rec.p, refracted}
			}

			return true
		}

		var n1, n2 float64

		if incoming.Dot(rec.normal) > 0 {
			n1 = rec.material.ior
			n2 = 1.0
		} else {
			n1 = 1.0
			n2 = rec.material.ior
		}

		reflectProbability = Fresnel(n1, n2, rec.normal, incoming) * (1 - phi)
		reflectProbability = (specularity + (1.0-specularity)*reflectProbability)

		if RandFloat(generator) <= reflectProbability {
			*scattered = Ray{rec.p, reflected}
			*attenuation = Color{reflectProbability, reflectProbability, reflectProbability}
		} else {
			target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
			*scattered = Ray{rec.p, target.Subtract(rec.p)}
		}

		return true
	}
	return false
}
