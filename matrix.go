package matrix

import (
	"gonum.org/v1/gonum/mat"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix mat.Dense

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (m Matrix) Dims() (r, c int) {
	md := mat.Dense(m)
	r, c = md.Dims()
	return
}

func (m Matrix) Split(i, j int) [][]Matrix {
	md := mat.Dense(m)
	r, c := m.Dims()
	matrices := make([][]Matrix, i)
	for k := range matrices {
		matrices[k] = make([]Matrix, j)
		rMin := k * (r / i)
		rMax := min((k+1)*(i/r), r)
		for l := range matrices[k] {
			cMin := l * (c / i)
			cMax := min((l+1)*(c/i), c)
			matrices[k][l] = mat.DenseCopyOf(md.Slice(rMin, rMax, cMin, cMax))
		}
	}
}

func (m Matrix) Add(a, b Matrix) {
	md := mat.Dense(m)
	ad := mat.Dense(a)
	bd := mat.Dense(b)
	md.Add(&ad, &bd)
}

func (m Matrix) Mul(a, b Matrix) {
	md := mat.Dense(m)
	ad := mat.Dense(a)
	bd := mat.Dense(b)
	md.Mul(&ad, &bd)
}

func (m Matrix) Diffusion(a Matrix) {

}
