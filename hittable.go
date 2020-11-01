package main

import (
	"math"
	"sort"
)

type HitRecord struct {
	u, v, t  float64
	uT, vT   float64
	p        Tuple
	normal   Tuple
	material Material
}

type HittableList struct {
	sphereBvh BVHSphere
	bvh       []*BVH
	atm       []Atmosphere
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

type BVHSphere struct {
	left, right *BVHSphere
	leaves      [2]LeafSphere
	bounds      AABB
	last        bool
	depth       int
}

type LeafSphere struct {
	bounds  AABB
	spheres []Sphere
}

func getBoundingBox(triangles []Triangle) AABB {
	xMin, xMax, yMin, yMax, zMin, zMax := -math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64

	var aabb AABB
	for _, triangle := range triangles {
		x1 := triangle.position.vertex0.x
		x2 := triangle.position.vertex1.x
		x3 := triangle.position.vertex2.x
		tempMin := max3(x1, x2, x3)
		tempMax := min3(x1, x2, x3)
		xMin = math.Max(xMin, tempMin)
		xMax = math.Min(xMax, tempMax)

		y1 := triangle.position.vertex0.y
		y2 := triangle.position.vertex1.y
		y3 := triangle.position.vertex2.y
		tempMin = max3(y1, y2, y3)
		tempMax = min3(y1, y2, y3)
		yMin = math.Max(yMin, tempMin)
		yMax = math.Min(yMax, tempMax)

		z1 := triangle.position.vertex0.z
		z2 := triangle.position.vertex1.z
		z3 := triangle.position.vertex2.z
		tempMin = max3(z1, z2, z3)
		tempMax = min3(z1, z2, z3)
		zMin = math.Max(zMin, tempMin)
		zMax = math.Min(zMax, tempMax)
	}

	aabb.min = Tuple{xMax, yMax, zMax, 0}
	aabb.max = Tuple{xMin, yMin, zMin, 0}

	return aabb
}

func getBoundingBoxSpheres(spheres []Sphere) AABB {
	xMin, xMax, yMin, yMax, zMin, zMax := -math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64

	var aabb AABB
	for _, sphere := range spheres {
		r := sphere.radius
		x1 := sphere.origin.x + r
		x2 := sphere.origin.x - r
		xMax = math.Min(xMax, x2)
		xMin = math.Max(xMin, x1)

		y1 := sphere.origin.y + r
		y2 := sphere.origin.y - r
		yMax = math.Min(yMax, y2)
		yMin = math.Max(yMin, y1)

		z1 := sphere.origin.z + r
		z2 := sphere.origin.z - r
		zMax = math.Min(zMax, z2)
		zMin = math.Max(zMin, z1)
	}

	aabb.min = Tuple{xMax, yMax, zMax, 0}
	aabb.max = Tuple{xMin, yMin, zMin, 0}

	return aabb
}

func getBVH(triangles []Triangle, depth, x int) *BVH {
	if x > 2 {
		x = 0
	}
	currentBox := getBoundingBox(triangles)
	lenX, lenY, lenZ := currentBox.sizeX(), currentBox.sizeY(), currentBox.sizeZ()
	maxLen := max3(lenX, lenY, lenZ)
	if maxLen == lenX {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.x < triangles[j].position.vertex0.x
		})
	} else if maxLen == lenY {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.y < triangles[j].position.vertex0.y
		})
	} else if maxLen == lenZ {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.z < triangles[j].position.vertex0.z
		})
	}
	x++
	size := len(triangles) / 2
	rightList := triangles[:size]
	leftList := triangles[size:]
	aabbLeft := getBoundingBox(leftList)
	aabbRight := getBoundingBox(rightList)
	if size <= limitTriangles {
		return &BVH{
			&BVH{}, &BVH{},
			[2]Leaf{
				Leaf{aabbLeft, leftList},
				Leaf{aabbRight, rightList},
			},
			getBoundingBox(triangles),
			true,
			depth,
		}
	}
	if depth > 0 {
		return &BVH{
			getBVH(leftList, depth-1, x), getBVH(rightList, depth-1, x),
			[2]Leaf{},
			getBoundingBox(triangles),
			false,
			depth,
		}
	}
	return &BVH{
		&BVH{}, &BVH{},
		[2]Leaf{
			Leaf{aabbLeft, leftList},
			Leaf{aabbRight, rightList},
		},
		getBoundingBox(triangles),
		true,
		depth,
	}
}

func getBVHSphere(spheres []Sphere, depth, x int) *BVHSphere {
	x++
	if x > 2 {
		x = 0
	}
	if x == 0 {
		sort.Slice(spheres[:], func(i, j int) bool {
			return spheres[i].origin.x < spheres[j].origin.x
		})
	} else if x == 1 {
		sort.Slice(spheres[:], func(i, j int) bool {
			return spheres[i].origin.y < spheres[j].origin.y
		})
	} else if x == 2 {
		sort.Slice(spheres[:], func(i, j int) bool {
			return spheres[i].origin.z < spheres[j].origin.z
		})
	}
	size := len(spheres) / 2
	rightList := spheres[:size]
	leftList := spheres[size:]
	aabbLeft := getBoundingBoxSpheres(leftList)
	aabbRight := getBoundingBoxSpheres(rightList)
	if size <= 1 {
		return &BVHSphere{
			&BVHSphere{}, &BVHSphere{},
			[2]LeafSphere{
				LeafSphere{aabbLeft, leftList},
				LeafSphere{aabbRight, rightList},
			},
			getBoundingBoxSpheres(spheres),
			true,
			depth,
		}
	}
	if depth > 0 {
		return &BVHSphere{
			getBVHSphere(leftList, depth-1, x), getBVHSphere(rightList, depth-1, x),
			[2]LeafSphere{},
			getBoundingBoxSpheres(spheres),
			false,
			depth,
		}
	}
	return &BVHSphere{
		&BVHSphere{}, &BVHSphere{},
		[2]LeafSphere{
			LeafSphere{aabbLeft, leftList},
			LeafSphere{aabbRight, rightList},
		},
		getBoundingBoxSpheres(spheres),
		true,
		depth,
	}
}

func hitBVH(tree *BVH, level int, r Ray, tMin, tMax float64) []Leaf {
	temp := []Leaf{}
	if tree == nil {
		return nil
	}
	if tree.last {
		if tree.bounds.hit(r, tMin, tMax) {
			if tree.leaves[0].bounds.hit(r, tMin, tMax) {
				temp = append(temp, tree.leaves[0])
			}
			if tree.leaves[1].bounds.hit(r, tMin, tMax) {
				temp = append(temp, tree.leaves[1])
			}
			return temp
		}
		return nil
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

func hitBVHSphere(tree *BVHSphere, level int, r Ray, tMin, tMax float64) [][2]LeafSphere {
	temp := [][2]LeafSphere{}
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
			temp = hitBVHSphere(tree.left, level-1, r, tMin, tMax)
		}
		if tree.right.bounds.hit(r, tMin, tMax) {
			tr := hitBVHSphere(tree.right, level-1, r, tMin, tMax)
			temp = append(temp, tr...)
		}
	}
	return temp
}

func (h *HittableList) hit(r Ray, tMin, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax

	spheres := hitBVHSphere(&h.sphereBvh, h.sphereBvh.depth, r, tMin, tMax)

	for i := 0; i < len(spheres); i++ {
		for j := 0; j < 2; j++ {
			if spheres[i][j].bounds.hit(r, tMin, tMax) {
				for k := 0; k < len(spheres[i][j].spheres); k++ {
					if spheres[i][j].spheres[k].hit(r, tMin, closestSoFar, &tempRec) {
						hitAnything = true
						closestSoFar = tempRec.t
						*rec = tempRec
					}
				}
			}
		}
	}

	for i := 0; i < len(h.bvh); i++ {
		tris := hitBVH(h.bvh[i], h.bvh[i].depth, r, tMin, tMax)
		for j := 0; j < len(tris); j++ {
			for k := 0; k < len(tris[j].triangles); k++ {
				if tris[j].triangles[k].hit(r, tMin, closestSoFar, &tempRec) {
					hitAnything = true
					closestSoFar = tempRec.t
					*rec = tempRec
				}
			}
		}
	}

	return hitAnything
}

func (s Sphere) uv(p Tuple) (float64, float64) {
	d := s.origin.Subtract(p).Normalize()
	u := 0.5 - (math.Atan2(d.z, d.x))/(2*math.Pi)
	v := 0.5 + (math.Asin(d.y))/(math.Pi)
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
			u, v := s.uv(*&rec.p)
			*&rec.u, *&rec.v = u, v
			*&rec.uT, *&rec.vT = u, v
			if len(s.material.albedo.normalTexture) > 0 {
				uvw := buildFromW(rec.normal)
				*&rec.normal = uvw.local(s.material.albedo.normal(*rec)).Normalize()
			}
			*&rec.material = s.material
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / (2.0 * a)
		if temp < tMax && temp > tMin {
			*&rec.t = temp
			*&rec.p = r.Position(rec.t)
			*&rec.normal = (rec.p.Subtract(s.origin)).DivScalar(s.radius).Normalize()
			u, v := s.uv(*&rec.p)
			*&rec.u, *&rec.v = u, v
			*&rec.uT, *&rec.vT = u, v
			if len(s.material.albedo.normalTexture) > 0 {
				uvw := buildFromW(rec.normal)
				*&rec.normal = uvw.local(s.material.albedo.normal(*rec)).Normalize()
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

		if tri.material.albedo.mode == TriangleImageUV || len(tri.material.albedo.normalTexture) > 0 {
			vt1 := tri.vtexture.vertex0
			vt2 := tri.vtexture.vertex1
			vt3 := tri.vtexture.vertex2
			x := vt2.MulScalar(u).Add(vt3.MulScalar(v)).Add(vt1.MulScalar(1 - u - v))
			*&rec.uT = x.x
			*&rec.vT = x.y
		}

		if tri.smooth {
			vn1 := tri.vnormals.vertex0
			vn2 := tri.vnormals.vertex1
			vn3 := tri.vnormals.vertex2
			*&rec.normal = vn2.MulScalar(u).Add(vn3.MulScalar(v)).Add(vn1.MulScalar(1 - u - v))
		} else {
			*&rec.normal = tri.normal
		}

		if len(tri.material.albedo.normalTexture) > 0 {
			uvw := buildFromW(rec.normal)
			*&rec.normal = uvw.local(tri.material.albedo.normal(*rec))
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

func (box *AABB) sizeX() float64 {
	return math.Abs(box.max.x - box.min.x)
}

func (box *AABB) sizeY() float64 {
	return math.Abs(box.max.y - box.min.y)
}

func (box *AABB) sizeZ() float64 {
	return math.Abs(box.max.z - box.min.z)
}
