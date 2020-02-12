package main

// Color struct holds three color values
type Color struct {
	r, g, b float64
}

// Add adds color values from other color values
func (c Color) Add(c1 Color) Color {
	return Color{c.r + c1.r, c.g + c1.g, c.b + c1.b}
}

// Subtract subtracts color values from other color values
func (c Color) Subtract(c1 Color) Color {
	return Color{c.r - c1.r, c.g - c1.g, c.b - c1.b}
}

// MulScalar multiplies color values by a scalar
func (c Color) MulScalar(s float64) Color {
	return Color{c.r * s, c.g * s, c.b * s}
}

// DivScalar divides color values by a scalar
func (c Color) DivScalar(s float64) Color {
	return Color{c.r / s, c.g / s, c.b / s}
}

// Mul returns Hadamard product of two colors
func (c Color) Mul(c1 Color) Color {
	return Color{c.r * c1.r, c.g * c1.g, c.b * c1.b}
}

func Hex(c int) Color {
	return Color{float64(c>>16&0xff) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255}
}
