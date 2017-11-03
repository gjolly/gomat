package matrix

import (
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

// Split the matrix m into i*j sub-matrices
func (m Matrix) Split(i, j int) [][]Matrix {
	r, c := m.Dims()
	matrices := make([][]Matrix, i)
	for k := range matrices {
		matrices[k] = make([]Matrix, j)
		rMin := k * (r / i)
		rMax := min((k+1)*(i/r), r)
		for l := range matrices[k] {
			cMin := l * (c / i)
			cMax := min((l+1)*(c/i), c)
			matrices[k][l] = Matrix{mat.DenseCopyOf(m.Slice(rMin, rMax, cMin, cMax))}
		}
	}
	return matrices
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
