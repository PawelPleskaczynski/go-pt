package main

import (
	"math"
)

// Mat is a struct wrapping 2D float64 array
type Mat struct {
	mat [][]float64
}

// MatShape returns shape of the matrix
func (mat Mat) MatShape() [2]int {
	return [2]int{len(mat.mat[0]), len(mat.mat)}
	//            ^~~~~width~~~~~  ^~~height~~~
}

// MatTranspose transposes matrix
func (mat Mat) MatTranspose() Mat {
	tempArray := make([][]float64, len(mat.mat))
	returnMat := Mat{tempArray}

	for i := 0; i < len(mat.mat); i++ {
		for j := 0; j < len(mat.mat[0]); j++ {
			tempArray[j] = append(tempArray[j], mat.mat[i][j])
		}
	}

	return returnMat
}

// MatMul multiplies two matrices
func (mat Mat) MatMul(mat1 Mat) Mat {
	tempArray := make([][]float64, len(mat.mat))
	returnMat := Mat{tempArray}

	for i := 0; i < len(mat.mat); i++ {
		tempArray[i] = make([]float64, len(mat1.mat[0]))
		for j := 0; j < len(mat1.mat[0]); j++ {
			for k := 0; k < len(mat1.mat); k++ {
				tempArray[i][j] += mat.mat[i][k] * mat1.mat[k][j]
			}
		}
	}

	return returnMat
}

// TupMul multiplies a matrix by a tuple
func (mat Mat) TupMul(tup Tuple) Tuple {
	tupMat := Mat{[][]float64{
		{tup.x},
		{tup.y},
		{tup.z},
		{1},
	}}

	tempArray := make([][]float64, len(mat.mat))

	for i := 0; i < len(mat.mat); i++ {
		tempArray[i] = make([]float64, len(tupMat.mat[0]))
		for j := 0; j < len(tupMat.mat[0]); j++ {
			for k := 0; k < len(mat.mat[0]); k++ {
				tempArray[i][j] += mat.mat[i][k] * tupMat.mat[k][j]
			}
		}
	}

	return Tuple{tempArray[0][0], tempArray[1][0], tempArray[2][0], tempArray[3][0]}
}

// ScalarMul multiplies a matrix by a scalar
func (mat Mat) ScalarMul(s float64) Mat {
	shape := mat.MatShape()

	for i := 0; i < shape[0]; i++ {
		for j := 0; j < shape[1]; j++ {
			mat.mat[i][j] = mat.mat[i][j] * s
		}
	}

	return mat
}

// GetIdentityMatrix returns identity matrix of given size
func GetIdentityMatrix(size int) Mat {
	tempArray := make([][]float64, size)
	returnMat := Mat{tempArray}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i == j {
				tempArray[j] = append(tempArray[j], 1)
			} else {
				tempArray[j] = append(tempArray[j], 0)
			}
		}
	}

	return returnMat
}

// Submatrix extracts a submatrix from a matrix, removing a row and a column
func (mat Mat) Submatrix(row, column int) Mat {
	tempArray := make([][]float64, len(mat.mat)-1)
	for i := range tempArray {
		tempArray[i] = make([]float64, len(mat.mat)-1)
	}
	originalShape := mat.MatShape()
	returnMat := Mat{tempArray}

	x, y := 0, 0

	for i := 0; i < originalShape[0]; i++ {
		if i != row {
			for j := 0; j < originalShape[1]; j++ {
				if j != column {
					if y >= originalShape[1] {
						break
					}
					tempArray[x][y] = mat.mat[i][j]
					y++
				}
			}
			x++
			y = 0
		}
		if x > originalShape[0] {
			break
		}
	}

	return returnMat
}

// Determinant returns determinant of the matrix
func (mat Mat) Determinant() float64 {
	width := mat.MatShape()[0]
	result := 0.0

	if width == 2 {
		result = mat.mat[0][0]*mat.mat[1][1] - mat.mat[1][0]*mat.mat[0][1]
	} else {
		for i := 0; i < width; i++ {
			result += mat.mat[0][i] * mat.Cofactor(0, i)
		}
	}

	return result
}

// Minor returns determinant of submatrix of 3x3 matrix
func (mat Mat) Minor(row, column int) float64 {
	matB := mat.Submatrix(row, column)
	return matB.Determinant()
}

// Cofactor returns cofactor of 3x3 matrix
func (mat Mat) Cofactor(row, column int) float64 {
	minor := mat.Minor(row, column)
	row++
	column++
	if row*column%2 == 0 {
		return -minor
	}
	return minor
}

// IsInvertible checks if matrix is invertible
func (mat Mat) IsInvertible() bool {
	return mat.Determinant() != 0
}

// Invert inverts the matrix
func (mat Mat) Invert() Mat {
	tempArray := make([][]float64, len(mat.mat))
	for i := range tempArray {
		tempArray[i] = make([]float64, len(mat.mat))
	}
	originalShape := mat.MatShape()
	returnMat := Mat{tempArray}

	for i := 0; i < originalShape[0]; i++ {
		for j := 0; j < originalShape[1]; j++ {
			c := mat.Cofactor(i, j)
			tempArray[i][j] = c / mat.Determinant()
		}
	}

	return returnMat
}

// TranslationMat returns a matrix for translation by x, y, z
func TranslationMat(x, y, z float64) []Mat {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][3], transformMat.mat[1][3], transformMat.mat[2][3] = x, y, z

	return []Mat{transformMat, transformMat.MatTranspose().Invert()}
}

// ScaleMat returns a matrix for scaling by x, y, z
func ScaleMat(x, y, z float64) []Mat {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[1][1], transformMat.mat[2][2] = x, y, z

	return []Mat{transformMat, transformMat.MatTranspose().Invert()}
}

// RotateXMat returns a matrix for rotating in x axis by angle
func RotateXMat(angle float64) []Mat {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[1][1], transformMat.mat[2][1], transformMat.mat[1][2], transformMat.mat[2][2] = math.Cos(angle), -math.Sin(angle), math.Sin(angle), math.Cos(angle)

	return []Mat{transformMat, transformMat.MatTranspose().Invert()}
}

// RotateYMat returns a matrix for rotating in y axis by angle
func RotateYMat(angle float64) []Mat {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[2][0], transformMat.mat[0][2], transformMat.mat[2][2] = math.Cos(angle), math.Sin(angle), -math.Sin(angle), math.Cos(angle)

	return []Mat{transformMat, transformMat.MatTranspose().Invert()}
}

// RotateZMat returns a matrix for rotating in z axis by angle
func RotateZMat(angle float64) []Mat {
	transformMat := GetIdentityMatrix(4)
	transformMat.mat[0][0], transformMat.mat[0][1], transformMat.mat[1][0], transformMat.mat[1][1] = math.Cos(angle), -math.Sin(angle), math.Sin(angle), math.Cos(angle)

	return []Mat{transformMat, transformMat.MatTranspose().Invert()}
}
