package matrix

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix struct {
	*mat.Dense
}

// New creates a new Matrix
func New(r, c int, data []float64) *Matrix {
	return &Matrix{mat.NewDense(r, c, data)}
}

// ToString returns the string representation of the matrix.
func (m *Matrix) ToString() string {
	return fmt.Sprintf("%v", mat.Formatted(m, mat.Squeeze()))
}

// Equal returns whether the matrices a and b have the same size and are element-wise equal
func Equal(m1, m2 *Matrix) bool {
	return mat.Equal(m1, m2)
}

func Add(m1, m2 *Matrix) *Matrix {
	r, c := m1.Dims()
	m := New(r, c, nil)
	m.Add(m1, m2)
	return m
}

func Sub(m1, m2 *Matrix) *Matrix {
	r, c := m1.Dims()
	m := New(r, c, nil)
	m.Sub(m1, m2)
	return m
}

func Mul(m1, m2 *Matrix) *Matrix {
	r, _ := m1.Dims()
	_, c := m2.Dims()
	m := New(r, c, nil)
	m.Mul(m1, m2)
	return m
}
