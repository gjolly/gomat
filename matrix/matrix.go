package matrix

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix mat.Dense

func (m *Matrix) toDense() *mat.Dense {
	mDense := mat.Dense(*m)
	return &mDense
}

// New creates a new Matrix
func New(r, c int, data []float64) *Matrix {
	m := Matrix(*mat.NewDense(r, c, data))
	return &m
}

// ToString returns the string representation of the matrix.
func (m *Matrix) ToString() string {
	return fmt.Sprintf("%v", mat.Formatted(m.toDense(), mat.Squeeze()))
}

// Equal returns whether the matrices a and b have the same size and are element-wise equal
func Equal(m1, m2 *Matrix) bool {
	return mat.Equal(m1.toDense(), m2.toDense())
}

func Add(m1, m2 *Matrix) *Matrix {
	r, c := m1.toDense().Dims()
	m := mat.NewDense(r, c, nil)
	m.Add(m1.toDense(), m2.toDense())
	res := Matrix(*m)
	return &res
}

func Sub(m1, m2 *Matrix) *Matrix {
	r, c := m1.toDense().Dims()
	m := mat.NewDense(r, c, nil)
	m.Sub(m1.toDense(), m2.toDense())
	res := Matrix(*m)
	return &res
}

func Mul(m1, m2 *Matrix) *Matrix {
	r, c := m1.toDense().Dims()
	m := mat.NewDense(r, c, nil)
	m.Mul(m1.toDense(), m2.toDense())
	res := Matrix(*m)
	return &res
}

func MaxDim(m Matrix) int {
	x, y := m.toDense().Dims()
	if x > y {
		return x
	}
	return y
}
