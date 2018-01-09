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
	r, _ := m1.toDense().Dims()
	_, c := m2.toDense().Dims()
	m := mat.NewDense(r, c, nil)
	m.Mul(m1.toDense(), m2.toDense())
	res := Matrix(*m)
	return &res
}

func (m *Matrix) Dims() (int, int) {
	r, c := m.toDense().Dims()
	return r, c
}

func (m *Matrix) Slice(rMin, rMax, cMin, cMax int) *Matrix {
	res := Matrix(*m.toDense().Slice(rMin, rMax, cMin, cMax).(*mat.Dense))
	return &res
}

func (m *Matrix) RawRowView(i int) []float64 {
	return m.toDense().RawRowView(i)
}

func (m Matrix) MaxDim() int {
	x, y := m.toDense().Dims()
	if x > y {
		return x
	}
	return y
}
