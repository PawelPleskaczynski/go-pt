package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	hsize   = 1500
	vsize   = 1500
	samples = 8192
	depth   = 8
)

func colorize(r Ray, world *HittableList, d int, generator rand.Rand) Color {
	rec := HitRecord{}
	if world.hit(r, Epsilon, math.MaxFloat64, &rec) {
		var attenuation Color
		var scattered Ray
		if d < depth && rec.material.Scatter(r, rec, &attenuation, &scattered, generator) {
			if rec.material.material == Emission {
				return rec.material.albedo
			} else {
				return attenuation.Mul(colorize(scattered, world, d+1, generator))
			}
		} else {
			return Color{0, 0, 0}
		}
	} else {
		// unit_direction := r.direction.Normalize()
		// t := 0.5 * (unit_direction.y + 1.0)
		// return Color{1.0, 1.0, 1.0}.MulScalar(1.0 - t).Add(Color{0.5, 0.7, 1.0}.MulScalar(t))
		return Color{0, 0, 0}
	}
}

func loadOBJ(file *os.File, list *[]Triangle, material Material, smooth bool) {
	vertices := []Tuple{}
	vertNormals := []Tuple{}
	faceVerts := []TrianglePosition{}
	faceNormals := []TrianglePosition{}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Fields(scanner.Text())
		if len(text) > 0 {
			if text[0] == "v" {
				x, _ := strconv.ParseFloat(text[1], 64)
				y, _ := strconv.ParseFloat(text[2], 64)
				z, _ := strconv.ParseFloat(text[3], 64)
				vertices = append(vertices, Tuple{
					x, y, z, 0,
				})
			} else if text[0] == "vn" {
				x, _ := strconv.ParseFloat(text[1], 64)
				y, _ := strconv.ParseFloat(text[2], 64)
				z, _ := strconv.ParseFloat(text[3], 64)
				vertNormals = append(vertNormals, Tuple{
					x, y, z, 0,
				})
			} else if text[0] == "f" {
				values1 := strings.Split(text[1], "/")
				values2 := strings.Split(text[2], "/")
				values3 := strings.Split(text[3], "/")

				v1, _ := strconv.Atoi(values1[0])
				v2, _ := strconv.Atoi(values2[0])
				v3, _ := strconv.Atoi(values3[0])

				vn1, _ := strconv.Atoi(values1[2])
				vn2, _ := strconv.Atoi(values2[2])
				vn3, _ := strconv.Atoi(values3[2])

				faceVerts = append(faceVerts, TrianglePosition{
					vertices[v1-1], vertices[v2-1], vertices[v3-1],
				})

				faceNormals = append(faceNormals, TrianglePosition{
					vertNormals[vn1-1], vertNormals[vn2-1], vertNormals[vn3-1],
				})
			}
		}
	}

	for i := 0; i < len(faceVerts); i++ {
		triangle := Triangle{
			faceVerts[i],
			faceNormals[i],
			material,
			Tuple{0, 0, 0, 0},
			smooth,
		}
		vertex0 := faceNormals[i].vertex0
		vertex1 := faceNormals[i].vertex1
		vertex2 := faceNormals[i].vertex2
		triangle.normal = (vertex0.Add(vertex1).Add(vertex2)).Normalize()
		*list = append(*list, triangle)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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

func getBVH(triangles []Triangle, depth, x int) *BVH {
	x++
	if x > 2 {
		x = 0
	}
	if x == 0 {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.x < triangles[j].position.vertex0.x
		})
	} else if x == 1 {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.y < triangles[j].position.vertex0.y
		})
	} else if x == 2 {
		sort.Slice(triangles[:], func(i, j int) bool {
			return triangles[i].position.vertex0.z < triangles[j].position.vertex0.z
		})
	}
	size := len(triangles) / 2
	rightList := triangles[:size]
	leftList := triangles[size:]
	aabbLeft := getBoundingBox(leftList)
	aabbRight := getBoundingBox(rightList)
	if size <= 1 {
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

func main() {
	listSpheres := []Sphere{}
	listTriangles := []Triangle{}

	cameraPosition := Tuple{-2, 0.8, 0, 0}
	cameraDirection := Tuple{0, 0.5, 0, 0}
	focusDistance := cameraDirection.Subtract(cameraPosition).Magnitude()
	camera := getCamera(cameraPosition, cameraDirection, Tuple{0, 1, 0, 0}, 80, float64(hsize)/float64(vsize), 0.0, focusDistance)

	file, err := os.Open("light.obj")
	if err != nil {
		log.Fatal(err)
	}

	loadOBJ(file, &listTriangles, Material{Emission, Color{5, 5, 5}, 0, 1.45, 0, false}, false)

	file, err = os.Open("dragon.obj")
	if err != nil {
		log.Fatal(err)
	}

	loadOBJ(file, &listTriangles, Material{Plastic, Color{0.482, 0.858, 0.52}, 0, 1.5, 0.5, false}, true)

	listSpheres = append(listSpheres, Sphere{
		Tuple{-0.75, 0.25, 0, 0}, 0.25,
		Material{Dielectric, Color{1, 1, 1}, 0, 1.45, 0.5, false},
	})

	listSpheres = append(listSpheres, Sphere{
		Tuple{-0.75, 0.25, -0.55, 0}, 0.25,
		Material{Plastic, Color{1, 1, 1}, 0, 1.45, 0.5, false},
	})

	log.Println("Building BVHs...")
	bvh := getBVH(listTriangles, 12, 0)
	log.Println("Built BVHs")

	// BOTTOM
	listSpheres = append(listSpheres, Sphere{
		Tuple{0, -10000, 0, 0}, 10000,
		Material{Lambertian, Color{1, 1, 1}, 0, 1.45, 0, false},
	})

	world := HittableList{listSpheres, *bvh}

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	buf := make([][]Color, cpus)

	for i := 0; i < cpus; i++ {
		buf[i] = make([]Color, vsize*hsize)
	}

	ch := make(chan int, cpus)

	canvas := make([]Color, vsize*hsize)

	start := time.Now()

	samplesCPU := samples / cpus
	remainder := 0

	if samples < cpus {
		cpus = samples
		samplesCPU = 1
	} else if samples%cpus != 0 {
		remainder = samples % cpus
	}

	doneSamples := 0

	log.Printf("Rendering %d triangles and %d spheres at %dx%d at %d samples on %d cores\n", len(listTriangles), len(listSpheres), hsize, vsize, samples, cpus)

	for i := 0; i < cpus; i++ {
		go func(i int) {
			for s := 0; s < samplesCPU; s++ {
				source := rand.NewSource(time.Now().UnixNano())
				generator := rand.New(source)
				sample := time.Now()
				for y := vsize - 1; y >= 0; y-- {
					for x := 0; x < hsize; x++ {
						col := Color{0, 0, 0}
						u := (float64(x) + RandFloat(*generator)) / float64(hsize)
						v := (float64(y) + RandFloat(*generator)) / float64(vsize)
						r := camera.getRay(u, v, *generator)

						col = colorize(r, &world, 0, *generator)

						buf[i][y*hsize+x] = buf[i][y*hsize+x].Add(col)
					}
				}

				doneSamples++
				sampleTime := time.Since(sample)
				fmt.Printf("\r%.2f%% (% 3d/% 3d) % 15s/sample, % 15s sample time, ETA: % 15s", float64(doneSamples)/float64(samples)*100, doneSamples, samples, sampleTime, sampleTime/(vsize*hsize), sampleTime*(time.Duration(samples)-time.Duration(doneSamples))/time.Duration(cpus))
			}
			ch <- 1
		}(i)
	}

	for i := 0; i < cpus; i++ {
		<-ch
	}
	close(ch)

	if remainder != 0 {
		ch = make(chan int, cpus)
		log.Printf("\nRendering additional %d samples...\n", remainder)
		for i := 0; i < remainder; i++ {
			go func(i int) {
				source := rand.NewSource(time.Now().UnixNano())
				generator := rand.New(source)
				sample := time.Now()
				for y := vsize - 1; y >= 0; y-- {
					for x := 0; x < hsize; x++ {
						col := Color{0, 0, 0}
						u := (float64(x) + RandFloat(*generator)) / float64(hsize)
						v := (float64(y) + RandFloat(*generator)) / float64(vsize)
						r := camera.getRay(u, v, *generator)

						col = colorize(r, &world, 0, *generator)

						buf[i][y*hsize+x] = buf[i][y*hsize+x].Add(col)
					}
				}

				doneSamples++
				sampleTime := time.Since(sample)
				fmt.Printf("\r%.2f%% (% 3d/% 3d) % 15s/sample, % 15s sample time, ETA: % 15s", float64(doneSamples)/float64(samples)*100, doneSamples, samples, sampleTime, sampleTime/(vsize*hsize), sampleTime*(time.Duration(samples)-time.Duration(doneSamples))/time.Duration(cpus))
				ch <- 1
			}(i)
		}

		for i := 0; i < remainder; i++ {
			<-ch
		}
		close(ch)
	}

	elapsed := time.Since(start)
	log.Printf("\nRendering took %s\n", elapsed)
	for i := 0; i < cpus; i++ {
		for y := 0; y < vsize; y++ {
			for x := 0; x < hsize; x++ {
				canvas[y*hsize+x] = canvas[y*hsize+x].Add(buf[i][y*hsize+x])
			}
		}
	}

	for y := 0; y < vsize; y++ {
		for x := 0; x < hsize; x++ {
			canvas[y*hsize+x] = canvas[y*hsize+x].DivScalar(float64(samples))
		}
	}

	fmt.Printf("Saving...\n")
	// filename := fmt.Sprintf("frame_%d.ppm", 0)
	filename := fmt.Sprintf("frame_%d", time.Now().UnixNano()/1e6)

	SaveImage(canvas, hsize, vsize, 255, filename, PNG, 16)
}
