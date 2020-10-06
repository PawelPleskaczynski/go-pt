package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadMaterial(file *os.File, name string, imageArray *[]ImageHash) Material {
	material := Material{material: Lambertian}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Fields(scanner.Text())
		if len(text) > 0 {
			if text[0] == "newmtl" && text[1] == name {
				for scanner.Scan() {
					text = strings.Fields(scanner.Text())
					if len(text) > 0 {
						if text[0] == "newmtl" {
							break
						}
						if text[0] == "Ke" {
							r, _ := strconv.ParseFloat(text[1], 64)
							g, _ := strconv.ParseFloat(text[2], 64)
							b, _ := strconv.ParseFloat(text[3], 64)
							if r > 0.0 || g > 0.0 || b > 0.0 {
								material.albedo = getConstant(Color{r, g, b})
								material.material = Emission
							}
						}
						if text[0] == "Kd" {
							r, _ := strconv.ParseFloat(text[1], 64)
							g, _ := strconv.ParseFloat(text[2], 64)
							b, _ := strconv.ParseFloat(text[3], 64)
							material.albedo = getConstant(Color{r, g, b})
						}
						if text[0] == "Ks" {
							r, _ := strconv.ParseFloat(text[1], 64)
							g, _ := strconv.ParseFloat(text[2], 64)
							b, _ := strconv.ParseFloat(text[3], 64)
							material.clearcoat = (r + g + b) / 3
						}
						if text[0] == "Ns" {
							roughness, _ := strconv.ParseFloat(text[1], 64)
							x1, _ := solveQuadratic(900, -1800, 900-roughness)
							if x1 > 1.0 {
								x1 = 1.0
							} else if x1 < 0.0 {
								x1 = 0.0
							}
							material.roughness = x1
							material.clearcoatRoughness = x1
						}
						if text[0] == "Ni" {
							ior, _ := strconv.ParseFloat(text[1], 64)
							material.ior = ior
						}
						if text[0] == "d" {
							transmission, _ := strconv.ParseFloat(text[1], 64)
							material.transmission = 1 - transmission
						}
						if text[0] == "Tr" {
							transmission, _ := strconv.ParseFloat(text[1], 64)
							material.transmission = transmission
						}
						if text[0] == "illum" {
							if material.material != Emission {
								mode, _ := strconv.ParseInt(text[1], 0, 0)
								switch mode {
								case 1:
									material.material = Lambertian
								case 2, 4, 6, 7, 9:
									material.material = BSDF
								case 3:
									material.material = BSDF
									material.metalicity = 1.0
								}
							}
						}
						if text[0] == "map_Kd" || text[0] == "map_Ke" {
							if fileExists(text[1]) {
								texture := getTexture(text[1], imageArray)
								material.albedo.diffuseTexture = texture
								material.albedo.mode = TriangleImageUV
							}
						}
						if text[0] == "map_Bump" || text[0] == "map_bump" || text[0] == "bump" {
							if fileExists(text[1]) {
								texture := getTexture(text[1], imageArray)
								material.albedo.normalTexture = texture
							}
						}
					}
				}
			}
		}
	}
	return material
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadOBJ(path string, list *[][]Triangle, transformationMatrix Mat, imageArray *[]ImageHash, materialArray *[]MaterialHash, material Material, smooth, overrideMaterial bool) {
	log.Printf("Loading 3D scene from %v file\n", path)
	vertices := []Tuple{}
	vertNormals := []Tuple{}
	vertTexture := []Tuple{}
	faceVerts := []TrianglePosition{}
	faceNormals := []TrianglePosition{}
	faceTexture := []TrianglePosition{}
	var materialFile *os.File

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	f := 0

	exists := false

	object := []Triangle{}

	defer file.Close()
	defer materialFile.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		materialFile.Seek(0, io.SeekStart)
		text := strings.Fields(scanner.Text())
		if len(text) > 0 {
			if text[0] == "o" {
				if len(object) > 0 {
					*list = append(*list, object)
					object = nil
				}
			}
			if !overrideMaterial {
				if text[0] == "mtllib" {
					if fileExists(text[1]) {
						log.Printf("Opening material library file %s", text[1])
						materialFile, _ = os.Open(text[1])
						exists = true
					} else {
						exists = false
					}
				}
				if text[0] == "usemtl" {
					if exists {
						strHash := hash(text[1])
						result := wasMaterialLoaded(strHash, *materialArray)
						if result == -1 {
							log.Printf("Loading new material: %s", text[1])
							material = loadMaterial(materialFile, text[1], imageArray)

							*materialArray = append(*materialArray, MaterialHash{
								material, strHash,
							})
						} else {
							material = (*materialArray)[result].material
						}
					}
				}
			}
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
			} else if text[0] == "vt" {
				u, _ := strconv.ParseFloat(text[1], 64)
				v, _ := strconv.ParseFloat(text[2], 64)
				vertTexture = append(vertTexture, Tuple{
					u, v, 1, 0,
				})
			} else if text[0] == "f" {
				vertCount := len(text) - 1
				for i := 0; i < vertCount-2; i++ {
					values1 := strings.Split(text[1], "/")
					values2 := strings.Split(text[i+2], "/")
					values3 := strings.Split(text[i+3], "/")

					var v1, v2, v3, vt1, vt2, vt3, vn1, vn2, vn3 int

					v1, _ = strconv.Atoi(values1[0])
					v2, _ = strconv.Atoi(values2[0])
					v3, _ = strconv.Atoi(values3[0])

					if values1[1] != "" && values2[1] != "" && values3[1] != "" {
						vt1, _ = strconv.Atoi(values1[1])
						vt2, _ = strconv.Atoi(values2[1])
						vt3, _ = strconv.Atoi(values3[1])

						if vt1 < 0 {
							vt1 = len(vertNormals) + vt1 + 1
						}
						if vt2 < 0 {
							vt2 = len(vertNormals) + vt2 + 1
						}
						if vt3 < 0 {
							vt3 = len(vertNormals) + vt3 + 1
						}

						faceTexture = append(faceTexture, TrianglePosition{
							vertTexture[vt1-1], vertTexture[vt2-1], vertTexture[vt3-1],
						})
					} else {
						faceTexture = append(faceTexture, TrianglePosition{})
					}

					if values1[2] != "" && values2[2] != "" && values3[2] != "" {
						vn1, _ = strconv.Atoi(values1[2])
						vn2, _ = strconv.Atoi(values2[2])
						vn3, _ = strconv.Atoi(values3[2])

						if vn1 < 0 {
							vn1 = len(vertNormals) + vn1 + 1
						}
						if vn2 < 0 {
							vn2 = len(vertNormals) + vn2 + 1
						}
						if vn3 < 0 {
							vn3 = len(vertNormals) + vn3 + 1
						}

						faceNormals = append(faceNormals, TrianglePosition{
							vertNormals[vn1-1], vertNormals[vn2-1], vertNormals[vn3-1],
						})
					} else {
						faceNormals = append(faceNormals, TrianglePosition{})
					}

					if v1 < 0 {
						v1 = len(vertices) + v1 + 1
					}
					if v2 < 0 {
						v2 = len(vertices) + v2 + 1
					}
					if v3 < 0 {
						v3 = len(vertices) + v3 + 1
					}

					faceVerts = append(faceVerts, TrianglePosition{
						vertices[v1-1], vertices[v2-1], vertices[v3-1],
					})

					triangle := Triangle{
						faceVerts[f],
						faceTexture[f],
						faceNormals[f],
						material,
						Tuple{0, 0, 0, 0},
						smooth,
					}
					vertex0 := faceNormals[f].vertex0
					vertex1 := faceNormals[f].vertex1
					vertex2 := faceNormals[f].vertex2
					triangle.normal = (vertex0.Add(vertex1).Add(vertex2)).Normalize()

					if len(transformationMatrix.mat) > 0 {
						triangle.normal = transformationMatrix.TupMul(triangle.normal)
						triangle.position.vertex0 = transformationMatrix.TupMul(triangle.position.vertex0)
						triangle.position.vertex1 = transformationMatrix.TupMul(triangle.position.vertex1)
						triangle.position.vertex2 = transformationMatrix.TupMul(triangle.position.vertex2)
						triangle.vnormals.vertex0 = transformationMatrix.TupMul(triangle.vnormals.vertex0)
						triangle.vnormals.vertex1 = transformationMatrix.TupMul(triangle.vnormals.vertex1)
						triangle.vnormals.vertex2 = transformationMatrix.TupMul(triangle.vnormals.vertex2)
						triangle.vtexture.vertex0 = transformationMatrix.TupMul(triangle.vtexture.vertex0)
						triangle.vtexture.vertex1 = transformationMatrix.TupMul(triangle.vtexture.vertex1)
						triangle.vtexture.vertex2 = transformationMatrix.TupMul(triangle.vtexture.vertex2)
					}

					object = append(object, triangle)
					f++
				}
			}
		}
	}

	if len(object) > 0 {
		*list = append(*list, object)
		object = nil
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
