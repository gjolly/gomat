package gomatcore

import (
	"math"

	"github.com/matei13/gomat/matrix"
	"gonum.org/v1/gonum/mat"
)

// SubMatrix represent a block of a matrix.
type SubMatrix struct {
	Mat *matrix.Matrix // Block
	Row int            // Row of the block in the original matrix
	Col int            // Column of the block in the original matrix
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// New creates a new Matrix
func New(r, c int, data []float64) *matrix.Matrix {
	return &matrix.Matrix{mat.NewDense(r, c, data)}
}

// Split the matrix m into sub-matrices of max size (n, n)
func Split(m *matrix.Matrix, n int) []*SubMatrix {
	r, c := m.Dims()
	nbRow := int(math.Ceil(float64(r) / float64(n)))
	nbCol := int(math.Ceil(float64(c) / float64(n)))
	matrices := make([]*SubMatrix, nbRow*nbCol)
	for i := 0; i < nbRow; i++ {
		rMin := i * n
		rMax := min((i+1)*n, r)
		for j := 0; j < nbCol; j++ {
			cMin := j * n
			cMax := min((j+1)*n, c)
			matrices[i*nbCol+j] = &SubMatrix{
				Mat: &matrix.Matrix{m.Slice(rMin, rMax, cMin, cMax).(*mat.Dense)},
				Row: i,
				Col: j,
			}
		}
	}
	return matrices
}

/*func mergedSize(matrices [][]*Matrix) (int, int) {
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
}*/

// Merge i*j sub-matrices into one matrix
func Merge(matrices []*SubMatrix, rM, cM, n int) *matrix.Matrix {
	//rM, cM := mergedSize(matrices)

	data := make([]float64, rM*cM)

	for _, sm := range matrices {
		r, c := sm.Mat.Dims()
		for k := 0; k < r; k++ {
			pos := (sm.Row*n+k)*cM + sm.Col*n
			copy(data[pos:pos+c], sm.Mat.RawRowView(k))
		}
	}

	m := matrix.New(rM, cM, data)

	return m
}
