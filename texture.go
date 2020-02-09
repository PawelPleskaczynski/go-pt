package main

import (
	"math"
)

const (
	Constant = iota
	Checkerboard
	CheckerboardUV
	ImageUV
)

type Texture struct {
	c                      []Color
	scaleX, scaleY, scaleZ float64
	mode                   int
	texture                [][]Color
}

func getConstant(c Color) Texture {
	return Texture{[]Color{c}, 0, 0, 0, Constant, nil}
}

func getCheckerboard(c1, c2 Color, scaleX, scaleY, scaleZ float64) Texture {
	return Texture{[]Color{c1, c2}, scaleX, scaleY, scaleZ, Checkerboard, nil}
}

func getCheckerboardUV(c1, c2 Color, scaleU, scaleV float64) Texture {
	return Texture{[]Color{c1, c2}, scaleU, scaleV, 0, CheckerboardUV, nil}
}

func getImageUV(texture [][]Color) Texture {
	return Texture{nil, 0, 0, 0, ImageUV, texture}
}

func (t Texture) color(u, v float64, p Tuple) Color {
	if t.mode == Constant {
		return t.c[0]
	} else if t.mode == Checkerboard {
		if (int(math.Floor(p.x/t.scaleX))+int(math.Floor(p.y/t.scaleY))+int(math.Floor(p.z/t.scaleZ)))%2 == 0 {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == CheckerboardUV {
		if (int(math.Floor(u/t.scaleX))+int(math.Floor(v/t.scaleY)))%2 == 0 {
			return t.c[0]
		}
		return t.c[1]
	} else if t.mode == ImageUV {
		nx := float64(len(t.texture))
		ny := float64(len(t.texture[0]))
		i := u * nx
		j := v*ny - 0.001
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
