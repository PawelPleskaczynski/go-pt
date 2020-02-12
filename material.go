package main

import (
	"math/rand"
)

const (
	Metal = iota
	Lambertian
	Dielectric
	Emission
	BSDF
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

func (m Material) Scatter(r Ray, rec HitRecord, attenuation *Color, scattered *Ray, generator rand.Rand) bool {
	if m.material == Lambertian {
		target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
		*scattered = Ray{rec.p, target.Subtract(rec.p)}
		*attenuation = m.albedo.color(rec.u, rec.v, rec.p)
		return true
	} else if m.material == Metal {
		reflected := (r.direction.Normalize()).Reflection(rec.normal)
		*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
		*attenuation = m.albedo.color(rec.u, rec.v, rec.p)
		return (scattered.direction.Dot(rec.normal) > 0)
	} else if m.material == Dielectric {
		var outwardNormal Tuple
		var refracted Tuple

		var niOverNt float64
		var reflectProbability float64
		var cosine float64

		*attenuation = m.albedo.color(rec.u, rec.v, rec.p)
		reflected := r.direction.Reflection(rec.normal)

		if r.direction.Dot(rec.normal) > 0 {
			outwardNormal = rec.normal.MulScalar(-1)
			niOverNt = m.ior
			cosine = m.ior * r.direction.Dot(rec.normal) / r.direction.Magnitude()
		} else {
			outwardNormal = rec.normal
			niOverNt = 1.0 / m.ior
			cosine = -(r.direction.Dot(rec.normal) / r.direction.Magnitude())
		}

		if r.direction.Refraction(outwardNormal, niOverNt, &refracted) {
			reflectProbability = Schlick(cosine, m.ior) + m.specularity/2
			if reflectProbability > 1.0 {
				reflectProbability = 1.0
			}
		} else {
			reflectProbability = 1.0
		}

		if RandFloat(generator) < reflectProbability {
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

		*attenuation = m.albedo.color(rec.u, rec.v, rec.p)
		reflected := r.direction.Reflection(rec.normal)

		if RandFloat(generator) < m.metalicity {
			*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
			return (scattered.direction.Dot(rec.normal) > 0)
		}

		if RandFloat(generator) < m.transmission {
			if r.direction.Dot(rec.normal) > 0 {
				outwardNormal = rec.normal.MulScalar(-1)
				niOverNt = m.ior
				cosine = m.ior * r.direction.Dot(rec.normal) / r.direction.Magnitude()
			} else {
				outwardNormal = rec.normal
				niOverNt = 1.0 / m.ior
				cosine = -(r.direction.Dot(rec.normal) / r.direction.Magnitude())
			}

			if r.direction.Refraction(outwardNormal, niOverNt, &refracted) {
				reflectProbability = Schlick(cosine, m.ior) + m.specularity/2
				if reflectProbability > 1.0 {
					reflectProbability = 1.0
				}
			} else {
				reflectProbability = 1.0
			}

			if RandFloat(generator) < reflectProbability {
				*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
				*attenuation = Color{reflectProbability, reflectProbability, reflectProbability}
			} else {
				*scattered = Ray{rec.p, refracted.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
			}

			return true
		}

		if r.direction.Dot(rec.normal) > 0 {
			outwardNormal = rec.normal.MulScalar(-1)
			niOverNt = m.ior
			cosine = 1 * r.direction.Dot(rec.normal) / r.direction.Magnitude()
		} else {
			outwardNormal = rec.normal
			niOverNt = 1.0 / m.ior
			cosine = -(r.direction.Dot(rec.normal) / r.direction.Magnitude())
		}

		if r.direction.Refraction(outwardNormal, niOverNt, &refracted) {
			reflectProbability = Schlick(cosine, m.ior) + m.specularity/2
			if reflectProbability > 1.0 {
				reflectProbability = 1.0
			}
		} else {
			reflectProbability = 1.0
		}

		if RandFloat(generator) < reflectProbability {
			*scattered = Ray{rec.p, reflected.Add(RandInUnitSphere(generator).MulScalar(m.roughness))}
			*attenuation = Color{reflectProbability, reflectProbability, reflectProbability}
		} else {
			target := rec.p.Add(rec.normal).Add(RandInUnitSphere(generator))
			*scattered = Ray{rec.p, target.Subtract(rec.p)}
		}

		return true
	}
	return false
}
