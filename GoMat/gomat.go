package gomat

import (
	"gonum.org/v1/gonum/mat"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix struct {
	*mat.Dense
}

// Add : Addition of two matrices
func Add(m1, m2 *Matrix) *Matrix {
	// TODO: Call deamon, wait for result and returns it
	return nil
}

// Sub : Substruction of two matrices
func Sub(m1, m2 *Matrix) *Matrix {
	// TODO: As above
	return nil
}

// Mult : Multiplication of two matrices
func Mult(m1, m2 *Matrix) *Matrix {
	// TODO: As above
	return nil
}
