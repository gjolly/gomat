package gomatcore

import (
	"math"

	"github.com/matei13/gomat/matrix"
)

// SubMatrix represent a block of a matrix.
type SubMatrix struct {
	Mat *matrix.Matrix // Block
	Row uint32            // Row of the block in the original matrix
	Col uint32            // Column of the block in the original matrix
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
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
				Mat: m.Slice(rMin, rMax, cMin, cMax),
				Row: uint32(i),
				Col: uint32(j),
			}
		}
	}
	return matrices
}

func (m SubMatrix) MaxDim() int {
	return m.Mat.MaxDim()
}

// Merge i*j sub-matrices into one matrix
func Merge(matrices []*SubMatrix, rM, cM, n int) *matrix.Matrix {
	//rM, cM := mergedSize(matrices)

	data := make([]float64, rM*cM)

	for _, sm := range matrices {
		r, c := sm.Mat.Dims()
		for k := 0; k < r; k++ {
			pos := (int(sm.Row)*n+k)*cM + int(sm.Col)*n
			copy(data[pos:pos+c], sm.Mat.RawRowView(k))
		}
	}

	m := matrix.New(rM, cM, data)

	return m
}
