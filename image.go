package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	PPM = iota
	PNG
)

// checks if there's an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SaveImage(canvas []Color, width, height, maxValue int, fileName string, extension int, depth int) {
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if canvas[y*width+x].r > 1.0 {
				canvas[y*width+x].r = 1.0
			}

			if canvas[y*width+x].g > 1.0 {
				canvas[y*width+x].g = 1.0
			}

			if canvas[y*width+x].b > 1.0 {
				canvas[y*width+x].b = 1.0
			}
		}
	}

	if extension == PPM {
		f, err := os.Create(fileName + ".ppm")
		check(err)

		w := bufio.NewWriter(f)
		defer w.Flush()

		_, err = fmt.Fprintf(w, "P3\n%d %d\n%d\n", width, height, maxValue)
		check(err)

		for y := height - 1; y >= 0; y-- {
			for x := 0; x < width; x++ {
				_, err := fmt.Fprintf(w, "%d %d %d ", int(math.Sqrt(canvas[y*width+x].r)*255), int(math.Sqrt(canvas[y*width+x].g)*255), int(math.Sqrt(canvas[y*width+x].b)*255))
				check(err)
			}
			_, err := fmt.Fprint(w, "\n")
			check(err)
		}
	} else if extension == PNG {
		f, err := os.Create(fileName + ".png")
		check(err)
		if depth == 8 {
			image := image.NewRGBA(image.Rect(0, 0, width, height))
			for y := height - 1; y >= 0; y-- {
				for x := 0; x < width; x++ {
					image.SetRGBA(x, height-1-y, color.RGBA{uint8(math.Sqrt(canvas[y*width+x].r) * 255.9), uint8(math.Sqrt(canvas[y*width+x].g) * 255.9), uint8(math.Sqrt(canvas[y*width+x].b) * 255.9), 255})
				}
			}
			png.Encode(f, image)
		} else {
			image := image.NewRGBA64(image.Rect(0, 0, width, height))
			for y := height - 1; y >= 0; y-- {
				for x := 0; x < width; x++ {
					image.SetRGBA64(x, height-1-y, color.RGBA64{uint16(math.Sqrt(canvas[y*width+x].r) * 65535.9), uint16(math.Sqrt(canvas[y*width+x].g) * 65535.9), uint16(math.Sqrt(canvas[y*width+x].b) * 65535.9), 65535})
				}
			}
			png.Encode(f, image)
		}
	}
}
