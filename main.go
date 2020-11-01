package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	_ "github.com/mdouchement/hdr/codec/rgbe"
)

const (
	hsize          = 512 * 3
	vsize          = 384 * 3
	samples        = 2048 * 8
	depth          = 8
	limitTriangles = 100
	preview        = false
	jitter         = true
)

func colorize(r Ray, world *HittableList, d int, generator rand.Rand, envMap Texture) Color {
	rec := HitRecord{}
	if world.hit(r, Epsilon, math.MaxFloat64, &rec) {
		var attenuation Color
		var scattered Ray
		if !preview {
			if d < depth && rec.material.Scatter(r, rec, &attenuation, &scattered, generator) {
				if rec.material.material == Emission {
					return rec.material.albedo.color(rec)
				} else {
					return attenuation.Mul(colorize(scattered, world, d+1, generator, envMap))
				}
			} else {
				return Color{0, 0, 0}
			}
		} else {
			if d < depth && rec.material.Scatter(r, rec, &attenuation, &scattered, generator) {
				if rec.material.metalicity > 0.0 {
					return rec.material.albedo.color(rec).Mul(colorize(scattered, world, d+1, generator, envMap))
				} else if rec.material.transmission > 0.0 {
					return rec.material.albedo.color(rec).Mul(colorize(scattered, world, d+1, generator, envMap))
				} else {
					shadeAmount := Tuple{0, 1, 0, 0}.Dot(rec.normal)
					shadowMin := 0.5
					return rec.material.albedo.color(rec).MulScalar(shadeAmount*(1-shadowMin) + shadowMin)
				}
			} else {
				return Color{0, 0, 0}
			}
		}
	} else {
		if len(world.atm) == 1 {
			d := r.direction.Normalize()
			rec.uT = 0.5 - (math.Atan2(d.z, d.x))/(2*math.Pi)*-1
			rec.vT = 0.5 + (math.Asin(d.y))/(math.Pi)*-1
			rec.p = d

			color := world.atm[0].ComputeIncidentLight(Tuple{0, world.atm[0].earthRadius + 1, 0, 0}, rec.p, 0, math.MaxFloat64)
			return Color{color.x, color.y, color.z}
		} else if envMap.mode == SphereImageUV {
			d := r.direction.Normalize()
			rec.uT = 0.5 - (math.Atan2(d.z, d.x))/(2*math.Pi)*-1
			rec.vT = 0.5 + (math.Asin(d.y))/(math.Pi)*-1
		}
		return envMap.color(rec)
	}
}

func main() {
	log.Println("Loading scene...")
	listSpheres := []Sphere{}
	listTriangles := [][]Triangle{}
	imageArray := []ImageHash{}
	materialArray := []MaterialHash{}
	transformationMatrix := GetIdentityMatrix(4)
	transformationMatrix = transformationMatrix.MatMul(RotateYMat(math.Pi)[0])

	averageFrameTime := time.Duration(0.0)
	averageSampleTime := time.Duration(0.0)
	numTris, done := 0, 0

	cameraPosition := Tuple{-5, 1.25, -5, 0}
	cameraDirection := Tuple{0, 0.8, 0, 0}

	focusDistance := cameraDirection.Subtract(cameraPosition).Magnitude()
	fLength := 40.0 // mm
	fNumber := 1024.0
	camera := getCamera(cameraPosition, cameraDirection, Tuple{0, 1, 0, 0}, fLength, float64(hsize)/float64(vsize), fNumber, focusDistance)

	atm := NewEarthAtmosphere(Tuple{1.0, 0.5, 0.0, 0})

	loadOBJ("scene.obj", &listTriangles, transformationMatrix, &imageArray, &materialArray, Material{}, true, false)

	listSpheres = append(listSpheres, Sphere{
		Tuple{-1 + 0.4, 0.2, -1 - 0.2, 0}, 0.2,
		getLambertian(getConstant(Color{1, 1, 1})),
	})

	listSpheres = append(listSpheres, Sphere{
		Tuple{-1, 0.2, -1, 0}, 0.2,
		getMetal(getConstant(Color{1, 1, 1}), 0.3, 0.0, 0.0),
	})

	listSpheres = append(listSpheres, Sphere{
		Tuple{-1 - 0.4, 0.2, -1 + 0.2, 0}, 0.2,
		getGlossy(getCheckerboardUV(Hex(0xffffff), Hex(0), 0.1, 0.2), 0, 1.0),
	})

	listSpheres = append(listSpheres, Sphere{
		Tuple{0, -100000 - Epsilon, 0, 0}, 100000,
		getLambertian(getCheckerboard(Color{0.5, 0.5, 0.5}, Color{0.2, 0.2, 0.2}, 0.5, 0.5, 0.5)),
	})

	bvh := []*BVH{}

	log.Println("Building BVHs...")
	for i := 0; i < len(listTriangles); i++ {
		numTris += len(listTriangles[i])
	}
	for i := 0; i < len(listTriangles); i++ {
		bvh = append(bvh, getBVH(listTriangles[i], 24, 0))
		done += len(listTriangles[i])
		fmt.Printf("\r%.2f%% (%d/%d triangles, %d/%d objects)", float64(done)/float64(numTris)*100, done, numTris, i+1, len(listTriangles))
	}
	println("")
	sphereBVH := getBVHSphere(listSpheres, 0, 0)
	log.Println("Built BVHs")

	world := HittableList{*sphereBVH, bvh, []Atmosphere{atm}}

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

	envMap := getConstant(Hex(0))
	// envMap := getConstant(Hex(0xffffff))
	// envMap := getImageUV(getTexture("interior.hdr", &imageArray))

	log.Printf("Rendering %d objects (%d triangles) and %d spheres at %dx%d at %d samples on %d cores\n", len(listTriangles), numTris, len(listSpheres), hsize, vsize, samples, cpus)

	for i := 0; i < cpus; i++ {
		go func(i int) {
			for s := 0; s < samplesCPU; s++ {
				source := rand.NewSource(time.Now().UnixNano())
				generator := rand.New(source)
				sample := time.Now()
				for y := vsize - 1; y >= 0; y-- {
					for x := 0; x < hsize; x++ {
						col := Color{0, 0, 0}
						u := float64(x)
						v := float64(y)
						if jitter {
							u += RandFloat(*generator)
							v += RandFloat(*generator)
						}
						u /= float64(hsize)
						v /= float64(vsize)
						r := camera.getRay(u, v, *generator)

						col = colorize(r, &world, 0, *generator, envMap)

						buf[i][y*hsize+x] = buf[i][y*hsize+x].Add(col)
					}
				}

				doneSamples++
				sampleTime := time.Since(sample)
				averageFrameTime += sampleTime
				averageSampleTime += sampleTime / (vsize * hsize)
				fmt.Printf("\r%.2f%% (% 3d/% 3d) % 15s/frame, % 15s sample time, ETA: % 15s", float64(doneSamples)/float64(samples)*100, doneSamples, samples, sampleTime, sampleTime/(vsize*hsize), sampleTime*(time.Duration(samples)-time.Duration(doneSamples))/time.Duration(cpus))
			}
			ch <- 1
		}(i)
	}

	for i := 0; i < cpus; i++ {
		<-ch
	}
	close(ch)

	if remainder != 0 {
		println()
		ch = make(chan int, remainder)
		log.Printf("Rendering additional %d samples...\n", remainder)
		for i := 0; i < remainder; i++ {
			go func(i int) {
				source := rand.NewSource(time.Now().UnixNano())
				generator := rand.New(source)
				sample := time.Now()
				for y := vsize - 1; y >= 0; y-- {
					for x := 0; x < hsize; x++ {
						col := Color{0, 0, 0}
						u := float64(x)
						v := float64(y)
						if jitter {
							u += RandFloat(*generator)
							v += RandFloat(*generator)
						}
						u /= float64(hsize)
						v /= float64(vsize)
						r := camera.getRay(u, v, *generator)

						col = colorize(r, &world, 0, *generator, envMap)

						buf[i][y*hsize+x] = buf[i][y*hsize+x].Add(col)
					}
				}

				doneSamples++
				sampleTime := time.Since(sample)
				fmt.Printf("\r%.2f%% (% 3d/% 3d) % 15s/frame, % 15s sample time, ETA: % 15s", float64(doneSamples)/float64(samples)*100, doneSamples, samples, sampleTime, sampleTime/(vsize*hsize), sampleTime*(time.Duration(samples)-time.Duration(doneSamples))/time.Duration(remainder))
				ch <- 1
			}(i)
		}

		for i := 0; i < remainder; i++ {
			<-ch
		}
		close(ch)
	}

	println()

	elapsed := time.Since(start)
	log.Printf("Rendering took %s\nAverage frame time: %s, average sample time: %s\n", elapsed, averageFrameTime/samples, averageSampleTime/samples)
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

	SaveImage(canvas, hsize, vsize, 255, filename, PNG, 16, true)
}
