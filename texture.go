package main

import (
	"math"
)

const (
	Constant = iota
	Checkerboard
	CheckerboardUV
	Grid
	GridUV
	SphereImageUV
	TriangleImageUV
)

type Texture struct {
	c                             []Color
	scaleX, scaleY, scaleZ, width float64
	mode                          int
	diffuseTexture                [][]Color
	normalTexture                 [][]Color
}

func getConstant(c Color) Texture {
	return Texture{[]Color{c}, 0, 0, 0, 0, Constant, nil, nil}
}

func getCheckerboard(c1, c2 Color, scaleX, scaleY, scaleZ float64) Texture {
	return Texture{[]Color{c1, c2}, scaleX, scaleY, scaleZ, 0, Checkerboard, nil, nil}
}

func getCheckerboardUV(c1, c2 Color, scaleU, scaleV float64) Texture {
	return Texture{[]Color{c1, c2}, scaleU, scaleV, 0, 0, CheckerboardUV, nil, nil}
}

func getGrid(c1, c2 Color, scaleX, scaleY, scaleZ, width float64) Texture {
	return Texture{[]Color{c1, c2}, scaleX, scaleY, scaleZ, width, Grid, nil, nil}
}

func getGridUV(c1, c2 Color, scaleU, scaleV, width float64) Texture {
	return Texture{[]Color{c1, c2}, scaleU, scaleV, 0, width, GridUV, nil, nil}
}

func getImageUV(texture [][]Color) Texture {
	return Texture{nil, 0, 0, 0, 0, SphereImageUV, texture, nil}
}

func getDiffNormalUV(diffuse, normal [][]Color) Texture {
	return Texture{nil, 0, 0, 0, 0, SphereImageUV, diffuse, normal}
}

func (t Texture) normal(rec HitRecord) Tuple {
	if t.mode == SphereImageUV {
		nx := float64(len(t.normalTexture))
		ny := float64(len(t.normalTexture[0]))
		i := rec.uT * nx
		j := rec.vT*ny - 0.001
		if i < 0 {
			i = 0
		}
		if j < 0 {
			j = 0
		}
		if i > nx-1 {
			i = nx - 1
		}
		if j > ny-1 {
			j = ny - 1
		}
		pixel := t.normalTexture[int(i)][int(j)]
		return Tuple{pixel.r, pixel.g, pixel.b, 1}.MulScalar(2).AddScalar(-1).Normalize()
	}
	nx := float64(len(t.normalTexture))
	ny := float64(len(t.normalTexture[0]))
	i := rec.uT * nx
	j := ny - rec.vT*ny - 0.001
	if i < 0 {
		i = 0
	}
	if j < 0 {
		j = 0
	}
	if i > nx-1 {
		i = nx - 1
	}
	if j > ny-1 {
		j = ny - 1
	}
	pixel := t.normalTexture[int(i)][int(j)]
	return Tuple{pixel.r, pixel.g, pixel.b, 1}.MulScalar(2).AddScalar(-1).Normalize()
}

func (t Texture) color(rec HitRecord) Color {
	if t.mode == Constant {
		return t.c[0]
	} else if t.mode == Checkerboard {
		if (int(math.Floor(rec.p.x/t.scaleX))+int(math.Floor(rec.p.y/t.scaleY))+int(math.Floor(rec.p.z/t.scaleZ)))%2 == 0 {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == CheckerboardUV {
		if (int(math.Floor(rec.uT/t.scaleX))+int(math.Floor(rec.vT/t.scaleY)))%2 == 0 {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == Grid {
		if (rec.p.x/t.scaleX-math.Floor(rec.p.x/t.scaleX)) < t.width || (rec.p.y/t.scaleY-math.Floor(rec.p.y/t.scaleY)) < t.width || (rec.p.z/t.scaleZ-math.Floor(rec.p.z/t.scaleZ)) < t.width {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == GridUV {
		if (rec.uT/t.scaleX-math.Floor(rec.uT/t.scaleX)) < t.width || (rec.vT/t.scaleY-math.Floor(rec.vT/t.scaleY)) < t.width {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == SphereImageUV {
		nx := float64(len(t.diffuseTexture))
		ny := float64(len(t.diffuseTexture[0]))
		i := rec.uT * nx
		j := rec.vT*ny - 0.001
		if i < 0 {
			i = 0
		}
		if j < 0 {
			j = 0
		}
		if i > nx-1 {
			i = nx - 1
		}
		if j > ny-1 {
			j = ny - 1
		}
		if math.IsNaN(i) || math.IsNaN(j) || math.IsNaN(nx) || math.IsNaN(ny) {
			return Color{}
		}
		return t.diffuseTexture[int(i)][int(j)]
	} else if t.mode == TriangleImageUV {
		nx := float64(len(t.diffuseTexture))
		ny := float64(len(t.diffuseTexture[0]))
		i := rec.uT * nx
		j := ny - rec.vT*ny - 0.001
		if i < 0 {
			i = 0
		}
		if j < 0 {
			j = 0
		}
		if i > nx-1 {
			i = nx - 1
		}
		if j > ny-1 {
			j = ny - 1
		}
		if math.IsNaN(i) || math.IsNaN(j) || math.IsNaN(nx) || math.IsNaN(ny) {
			return Color{}
		}
		return t.diffuseTexture[int(i)][int(j)]
	}
	return Color{}
}
