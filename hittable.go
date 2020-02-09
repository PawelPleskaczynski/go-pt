package main

import (
	"math"
)

type HitRecord struct {
	u, v, t  float64
	p        Tuple
	normal   Tuple
	material Material
}

type HittableList struct {
	sphereHits []Sphere
	bvh        BVH
}

type AABB struct {
	min, max Tuple
}

type BVH struct {
	left, right *BVH
	leaves      [2]Leaf
	bounds      AABB
	last        bool
	depth       int
}

type Leaf struct {
	bounds    AABB
	triangles []Triangle
}

func hitBVH(tree *BVH, level int, r Ray, tMin, tMax float64) [][2]Leaf {
	temp := [][2]Leaf{}
	if tree == nil {
		return nil
	}
	if tree.last {
		if tree.bounds.hit(r, tMin, tMax) {
			return append(temp, tree.leaves)
		}
		return temp
	} else if level > 0 {
		if tree.left.bounds.hit(r, tMin, tMax) {
			temp = hitBVH(tree.left, level-1, r, tMin, tMax)
		}
		if tree.right.bounds.hit(r, tMin, tMax) {
			tr := hitBVH(tree.right, level-1, r, tMin, tMax)
			temp = append(temp, tr...)
		}
	}
	return temp
}

func (h *HittableList) hit(r Ray, tMin, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax
	for i := 0; i < len(h.sphereHits); i++ {
		if h.sphereHits[i].hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}

	current := &h.bvh

	tris := hitBVH(current, current.depth, r, tMin, tMax)

	for i := 0; i < len(tris); i++ {
		for j := 0; j < 2; j++ {
			if tris[i][j].bounds.hit(r, tMin, tMax) {
				for k := 0; k < len(tris[i][j].triangles); k++ {
					if tris[i][j].triangles[k].hit(r, tMin, closestSoFar, &tempRec) {
						hitAnything = true
						closestSoFar = tempRec.t
						*rec = tempRec
					}
				}
			}
		}
	}
	return hitAnything
}

// func (s Sphere) uv(p Tuple) (float64, float64) {
// 	phi := math.Atan2(p.z, p.x)
// 	theta := math.Asin(p.y)
// 	fmt.Printf("%v\n", theta)
// 	u := 1 - (phi+math.Pi)/(2*math.Pi)
// 	v := (theta + math.Pi/2) / math.Pi
// 	return u, v
// }

func (s Sphere) uv(p Tuple) (float64, float64) {
	d := s.origin.Subtract(p).Normalize()
	u := 0.5 - (math.Atan2(d.z, d.x))/(2*math.Pi)
	v := 0.5 + (math.Asin(d.y))/(math.Pi)
	// fmt.Printf("%v %v\n", u, v)
	return u, v
}

func (s *Sphere) hit(r Ray, tMin, tMax float64, rec *HitRecord) bool {
	oc := r.origin.Subtract(s.origin)
	a := r.direction.Dot(r.direction)
	b := 2.0 * oc.Dot(r.direction)
	c := oc.Dot(oc) - s.radius*s.radius
	discriminant := b*b - 4*a*c
	if discriminant > 0.0 {
		temp := (-b - math.Sqrt(discriminant)) / (2.0 * a)
		if temp < tMax && temp > tMin {
			*&rec.t = temp
			*&rec.p = r.Position(rec.t)
			*&rec.normal = (rec.p.Subtract(s.origin)).DivScalar(s.radius).Normalize()
			if s.material.albedo.mode == CheckerboardUV || s.material.albedo.mode == ImageUV {
				*&rec.u, *&rec.v = s.uv(*&rec.p)
			}
			*&rec.material = s.material
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / (2.0 * a)
		if temp < tMax && temp > tMin {
			*&rec.t = temp
			*&rec.p = r.Position(rec.t)
			*&rec.normal = (rec.p.Subtract(s.origin)).DivScalar(s.radius).Normalize()
			if s.material.albedo.mode == CheckerboardUV || s.material.albedo.mode == ImageUV {
				*&rec.u, *&rec.v = s.uv(*&rec.p)
			}
			*&rec.material = s.material
			return true
		}
	}
	return false
}

func (tri *Triangle) hit(r Ray, tMin, tMax float64, rec *HitRecord) bool {
	vertex0 := tri.position.vertex0
	vertex1 := tri.position.vertex1
	vertex2 := tri.position.vertex2
	edge1 := vertex1.Subtract(vertex0)
	edge2 := vertex2.Subtract(vertex0)
	h := r.direction.Cross(edge2)
	a := edge1.Dot(h)
	if a > -Epsilon && a < Epsilon {
		return false
	}
	f := 1.0 / a
	s := r.origin.Subtract(vertex0)
	u := f * s.Dot(h)
	if u < 0.0 || u > 1.0 {
		return false
	}
	q := s.Cross(edge1)
	v := f * r.direction.Dot(q)
	if v < 0.0 || u+v > 1.0 {
		return false
	}
	t := f * edge2.Dot(q)
	if t < tMax && t > tMin {
		*&rec.p = r.origin.Add(r.direction.MulScalar(t))
		*&rec.t = t
		*&rec.material = tri.material
		*&rec.u = u
		*&rec.v = v
		if tri.smooth {
			vn1 := tri.vnormals.vertex0
			vn2 := tri.vnormals.vertex1
			vn3 := tri.vnormals.vertex2
			*&rec.normal = vn2.MulScalar(u).Add(vn3.MulScalar(v)).Add(vn1.MulScalar(1 - u - v))
		} else {
			*&rec.normal = tri.normal
		}
		return true
	}
	return false
}

func (box *AABB) hit(r Ray, tMin, tMax float64) bool {
	dirFrac := Tuple{
		1.0 / r.direction.x,
		1.0 / r.direction.y,
		1.0 / r.direction.z,
		1,
	}

	t1 := (box.min.x - r.origin.x) * dirFrac.x
	t2 := (box.max.x - r.origin.x) * dirFrac.x
	t3 := (box.min.y - r.origin.y) * dirFrac.y
	t4 := (box.max.y - r.origin.y) * dirFrac.y
	t5 := (box.min.z - r.origin.z) * dirFrac.z
	t6 := (box.max.z - r.origin.z) * dirFrac.z

	tMin = math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tMax = math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	if tMax < 0 {
		return false
	}

	if tMin > tMax {
		return false
	}

	return true
}
