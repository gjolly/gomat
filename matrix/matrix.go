package matrix

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix struct {
	*mat.Dense
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// New creates a new Matrix
func New(r, c int, data []float64) *Matrix {
	m := Matrix{mat.NewDense(r, c, data)}
	return &m
}

// Split the matrix m into i*j sub-matrices (i sub-columns, j sub-rows)
func (m *Matrix) Split(i, j int) [][]*Matrix {
	r, c := m.Dims()
	rSize := int(math.Ceil(float64(r) / float64(i)))
	cSize := int(math.Ceil(float64(c) / float64(j)))
	matrices := make([][]*Matrix, i)
	for k := range matrices {
		matrices[k] = make([]*Matrix, j)
		rMin := k * rSize
		rMax := min((k+1)*rSize, r)
		for l := range matrices[k] {
			cMin := l * cSize
			cMax := min((l+1)*cSize, c)
			matrices[k][l] = &Matrix{m.Slice(rMin, rMax, cMin, cMax).(*mat.Dense)}
		}
	}
	return matrices
}

func mergedSize(matrices [][]*Matrix) (int, int) {
	rM, cM := 0, 0
	length := len(matrices[0])
	for i := range matrices {
		if len(matrices[i]) != length {
			panic(mat.ErrShape)
		}
		r, c := 0, 0
		for j := range matrices[i] {
			rIJ, cIJ := matrices[i][j].Dims()
			if (j != 0) && (r != rIJ) {
				panic(mat.ErrShape)
			}
			r = rIJ
			c += cIJ
		}
		if (i != 0) && (c != cM) {
			panic(mat.ErrShape)
		}
		rM += r
		cM = c
	}
	return rM, cM
}

// Merge i*j sub-matrices into one matrix
func Merge(matrices [][]*Matrix) *Matrix {
	rM, cM := mergedSize(matrices)

	data := make([]float64, rM*cM)

	pos := 0
	for i := range matrices {
		r, _ := matrices[i][0].Dims()
		for k := 0; k < r; k++ {
			for j := range matrices[i] {
				_, c := matrices[i][j].Dims()
				copy(data[pos:pos+c], matrices[i][j].RawRowView(k))
				pos += c
			}
		}
	}

	m := New(rM, cM, data)

	return m
}

/*func (m Matrix) Diffusion(it int) Matrix {
	r, c := m.Dims()
	input := Matrix{mat.NewDense(r, c, make([]float64, r*c))}
	input.Copy(&m)
	output := Matrix{mat.NewDense(r, c, make([]float64, r*c))}
	for t := 0; t < it; t++ {
		for i := 1; i < c-1; i++ {
			for j := 1; j < r-1; j++ {
				cc := m.At(i-1, j-1) + m.At(i, j-1) + m.At(i+1, j-1)
				cc += m.At(i-1, j) + m.At(i, j) + m.At(i+1, j)
				cc += m.At(i-1, j+1) + m.At(i, j+1) + m.At(i+1, j+1)
				output.Set(i, j, cc/9)
			}
		}
		// Swap output and input
		tmp := input
		input = output
		output = tmp
	}
	return input
}*/

// ToString returns the string representation of the matrix.
func (m *Matrix) ToString() string {
	return fmt.Sprintf("%v", mat.Formatted(m, mat.Squeeze()))
}

// Equal returns whether the matrices a and b have the same size and are element-wise equal
func Equal(m1, m2 *Matrix) bool {
	return mat.Equal(m1, m2)
}
