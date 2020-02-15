package main

import (
	"math"
)

// Tuple struct holds three coordinate values and w value
type Tuple struct {
	x, y, z, w float64
	//		 ^ 0 = vector
	//		 | 1 = point
}

// Equals checks if two tuples are equal
func (v Tuple) Equals(u Tuple) bool {
	if Equal(v.x, u.x) && Equal(v.y, u.y) && Equal(v.z, u.z) && Equal(v.w, u.w) {
		return true
	}
	return false
}

// Add adds two tuples together
func (v Tuple) Add(u Tuple) Tuple {
	return Tuple{v.x + u.x, v.y + u.y, v.z + u.z, v.w + u.w}
}

// Subtract subtracts two tuples
func (v Tuple) Subtract(u Tuple) Tuple {
	return Tuple{v.x - u.x, v.y - u.y, v.z - u.z, v.w - u.w}
}

// AddScalar adds a scalar to values
func (v Tuple) AddScalar(s float64) Tuple {
	return Tuple{v.x + s, v.y + s, v.z + s, v.w}
}

// Negate negates a tuple
func (v Tuple) Negate() Tuple {
	return Tuple{-v.x, -v.y, -v.z, -v.w}
}

// MulScalar multiplies tuple by a scalar
func (v Tuple) MulScalar(s float64) Tuple {
	return Tuple{v.x * s, v.y * s, v.z * s, v.w * s}
}

// DivScalar divides tuple by a scalar
func (v Tuple) DivScalar(s float64) Tuple {
	return Tuple{v.x / s, v.y / s, v.z / s, v.w / s}
}

// Magnitude returns magnitude of the tuple
func (v Tuple) Magnitude() float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y) + (v.z * v.z))
}

// Normalize normalizes the tuple
func (v Tuple) Normalize() Tuple {
	magnitude := v.Magnitude()

	return Tuple{
		v.x / magnitude,
		v.y / magnitude,
		v.z / magnitude,
		v.w / magnitude,
	}
}

// Dot returns dot product of two vectors
func (v Tuple) Dot(u Tuple) float64 {
	return v.x*u.x +
		v.y*u.y +
		v.z*u.z
}

// Cross returns cross product of two vectors
func (v Tuple) Cross(u Tuple) Tuple {
	return Tuple{
		v.y*u.z - v.z*u.y,
		v.z*u.x - v.x*u.z,
		v.x*u.y - v.y*u.x,
		v.w,
	}
}

// functions for transformations

// Translate translates tuple by x, y, z values
func (v Tuple) Translate(x, y, z float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][3], transformMat.mat[1][3], transformMat.mat[2][3] = x, y, z

	return transformMat.TupMul(v)
}

// Scale scales tuple by x, y, z values
func (v Tuple) Scale(x, y, z float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[1][1], transformMat.mat[2][2] = x, y, z

	return transformMat.TupMul(v)
}

// RotateX rotates tuple around X axis
func (v Tuple) RotateX(angle float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[1][1], transformMat.mat[2][1], transformMat.mat[2][1], transformMat.mat[2][2] = math.Cos(angle), -math.Sin(angle), math.Sin(angle), math.Cos(angle)

	return transformMat.TupMul(v)
}

// RotateY rotates tuple around Y axis
func (v Tuple) RotateY(angle float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[2][0], transformMat.mat[2][0], transformMat.mat[2][2] = math.Cos(angle), math.Sin(angle), -math.Sin(angle), math.Cos(angle)

	return transformMat.TupMul(v)
}

// RotateZ rotates tuple around Z axis
func (v Tuple) RotateZ(angle float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[0][1], transformMat.mat[1][0], transformMat.mat[1][1] = math.Cos(angle), -math.Sin(angle), math.Sin(angle), math.Cos(angle)

	return transformMat.TupMul(v)
}

// Shear shears tuple
func (v Tuple) Shear(xy, xz, yx, yz, zx, zy float64) Tuple {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][1], transformMat.mat[0][2], transformMat.mat[1][0], transformMat.mat[1][2], transformMat.mat[2][0], transformMat.mat[2][1] = xy, xz, yx, yz, zx, zy

	return transformMat.TupMul(v)
}

// Reflection returns reflection vector
func (v Tuple) Reflection(n Tuple) Tuple {
	return v.Subtract(n.MulScalar(2).MulScalar(v.Dot(n)))
}

func (v Tuple) Refraction(n Tuple, niOverNt float64, refracted *Tuple) bool {
	dt := -n.Dot(v)
	discriminant := niOverNt * niOverNt * (1.0 - dt*dt)
	if discriminant > 1 {
		return false
	}
	*refracted = v.MulScalar(niOverNt).Add(n.MulScalar(niOverNt*dt - math.Sqrt(1-discriminant)))
	return true
}

func Schlick(cos, ior float64) float64 {
	r0 := (1.0 - ior) / (1.0 + ior)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1-cos, 5)
}
