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
	texture                       [][]Color
}

func getConstant(c Color) Texture {
	return Texture{[]Color{c}, 0, 0, 0, 0, Constant, nil}
}

func getCheckerboard(c1, c2 Color, scaleX, scaleY, scaleZ float64) Texture {
	return Texture{[]Color{c1, c2}, scaleX, scaleY, scaleZ, 0, Checkerboard, nil}
}

func getCheckerboardUV(c1, c2 Color, scaleU, scaleV float64) Texture {
	return Texture{[]Color{c1, c2}, scaleU, scaleV, 0, 0, CheckerboardUV, nil}
}

func getGrid(c1, c2 Color, scaleX, scaleY, scaleZ, width float64) Texture {
	return Texture{[]Color{c1, c2}, scaleX, scaleY, scaleZ, width, Grid, nil}
}

func getGridUV(c1, c2 Color, scaleU, scaleV, width float64) Texture {
	return Texture{[]Color{c1, c2}, scaleU, scaleV, 0, width, GridUV, nil}
}

func getImageUV(texture [][]Color) Texture {
	return Texture{nil, 0, 0, 0, 0, SphereImageUV, texture}
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
		if (int(math.Floor(rec.u/t.scaleX))+int(math.Floor(rec.v/t.scaleY)))%2 == 0 {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == Grid {
		if (rec.p.x/t.scaleX-math.Floor(rec.p.x/t.scaleX)) < t.width || (rec.p.y/t.scaleY-math.Floor(rec.p.y/t.scaleY)) < t.width || (rec.p.z/t.scaleZ-math.Floor(rec.p.z/t.scaleZ)) < t.width {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == GridUV {
		if (rec.u/t.scaleX-math.Floor(rec.u/t.scaleX)) < t.width || (rec.v/t.scaleY-math.Floor(rec.v/t.scaleY)) < t.width {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == SphereImageUV {
		nx := float64(len(t.texture))
		ny := float64(len(t.texture[0]))
		i := rec.u * nx
		j := rec.v*ny - 0.001
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
		return t.texture[int(i)][int(j)]
	} else if t.mode == TriangleImageUV {
		nx := float64(len(t.texture))
		ny := float64(len(t.texture[0]))
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
		return t.texture[int(i)][int(j)]
	}
	return Color{}
}
