package matrix_test

import (
	"testing"

	"github.com/matei13/gomat/matrix"
)

var arithmeticTest = []struct {
	m1  *matrix.Matrix
	m2  *matrix.Matrix
	m2T *matrix.Matrix
	sum *matrix.Matrix // m1 + m2
	mul *matrix.Matrix // m1 * m2T
}{
	{
		matrix.New(3, 2, []float64{0, 1, 2, 3, 4, 5}),
		matrix.New(3, 2, []float64{1, 2, 4, 6, 8, 10}),
		matrix.New(2, 3, []float64{1, 4, 8, 2, 6, 10}),
		matrix.New(3, 2, []float64{1, 3, 6, 9, 12, 15}),
		matrix.New(3, 3, []float64{2, 6, 10, 8, 26, 46, 14, 46, 82}),
	},
}

func TestAdd(t *testing.T) {
	for _, tt := range arithmeticTest {
		sum := matrix.AddMatrix(tt.m1, tt.m2)
		if !matrix.Equal(sum, tt.sum) {
			t.Errorf("Add: expected\n%v\n, actual \n%v\n", tt.sum.ToString(), sum.ToString())
		}
	}
}

func TestSub(t *testing.T) {
	for _, tt := range arithmeticTest {
		sub := matrix.SubMatrix(tt.sum, tt.m1)
		if !matrix.Equal(sub, tt.m2) {
			t.Errorf("Sub: expected\n%v\n, actual \n%v\n", tt.m2.ToString(), sub.ToString())
		}
	}
}

func TestMul(t *testing.T) {
	for _, tt := range arithmeticTest {
		mul := matrix.MulMatrix(tt.m1, tt.m2T)
		if !matrix.Equal(mul, tt.mul) {
			t.Errorf("Mul: expected\n%v\n, actual \n%v\n", tt.mul.ToString(), mul.ToString())
		}
	}
}
