package main

import (
	"crypto/sha1"
	"encoding/hex"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

type ImageHash struct {
	image [][]Color
	hash  string
}

func wasImageLoaded(hash string, table []ImageHash) int {
	for i := 0; i < len(table); i++ {
		if table[i].hash == hash {
			return i
		}
	}
	return -1
}

type MaterialHash struct {
	material Material
	hash     string
}

func wasMaterialLoaded(hash string, table []MaterialHash) int {
	for i := 0; i < len(table); i++ {
		if table[i].hash == hash {
			return i
		}
	}
	return -1
}

func hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha)
}

func solveQuadratic(a, b, c float64) (float64, float64) {
	discriminant := b*b - 4*a*c
	if discriminant > 0 {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a), (-b + math.Sqrt(discriminant)) / (2.0 * a)
	} else if discriminant == 0 {
		return -b / (2.0 * a), -b / (2.0 * a)
	}

	return math.MaxFloat64, math.MaxFloat64
}

func min3(a, b, c float64) float64 {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max3(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

func loadImage(path string) image.Image {
	log.Printf("Loading image: %s...", path)
	var texture image.Image
	if fileExists(path) {
		textureFile, _ := os.Open(path)
		if strings.HasSuffix(strings.ToLower(path), "png") {
			texture, _ = png.Decode(textureFile)
		} else if strings.HasSuffix(strings.ToLower(path), "jpg") || strings.HasSuffix(strings.ToLower(path), "jpeg") {
			texture, _ = jpeg.Decode(textureFile)
		} else if strings.HasSuffix(strings.ToLower(path), "hdr") {
			texture, _, err := image.Decode(textureFile)
			check(err)
			return texture
		}
	}
	return texture
}

func loadTexture(texture image.Image) [][]Color {
	width := texture.Bounds().Dx()
	height := texture.Bounds().Dy()
	array := make([][]Color, width)
	for i := 0; i < width; i++ {
		array[i] = make([]Color, height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := texture.At(x, y).RGBA()
			array[x][y] = Color{float64(r>>8) / 255, float64(g>>8) / 255, float64(b>>8) / 255}
		}
	}

	return array
}

func getTexture(path string, imageArray *[]ImageHash) [][]Color {
	var texture [][]Color
	strHash := hash(path)
	result := wasImageLoaded(strHash, *imageArray)
	if result == -1 {
		texture = loadTexture(loadImage(path))
		*imageArray = append(*imageArray, ImageHash{
			texture, strHash,
		})
	} else {
		texture = (*imageArray)[result].image
	}
	return texture
}
