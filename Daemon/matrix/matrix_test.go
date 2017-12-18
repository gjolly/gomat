package matrix_test

import (
	"math/rand"
	"testing"
	"github.com/matei13/gomat/Daemon/matrix"
)

func createMatrix(r, c int, fct func(int) float64) *matrix.Matrix {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = fct(i)
	}
	return matrix.New(r, c, data)
}

var splitMergeTest = []struct {
	m     *matrix.Matrix     // input matrix
	i     int                // nb of sub-columns
	j     int                // nb of sub-rows
	split [][]*matrix.Matrix // expected result
}{
	{
		createMatrix(5, 3, func(i int) float64 { return float64(i) }), 3, 2, [][]*matrix.Matrix{
		[]*matrix.Matrix{
			matrix.New(3, 2, []float64{0, 1, 3, 4, 6, 7}),
			matrix.New(3, 1, []float64{2, 5, 8}),
		},
		[]*matrix.Matrix{
			matrix.New(2, 2, []float64{9, 10, 12, 13}),
			matrix.New(2, 1, []float64{11, 14}),
		},
	},
	},
}

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

func TestSplit(t *testing.T) {
	for _, tt := range splitMergeTest {
		split := (tt.m).Split(tt.i, tt.j)
		for i := range split {
			for j := range split[i] {
				if !matrix.Equal(split[i][j], tt.split[i][j]) {
					t.Errorf("Split(%d,%d): expected\n%v\n, actual \n%v\n", i, j, tt.split[i][j].ToString(), split[i][j].ToString())
				}
			}
		}
	}
}

func TestMerge(t *testing.T) {
	for _, tt := range splitMergeTest {
		merged := matrix.Merge(tt.split)
		if !matrix.Equal(merged, tt.m) {
			t.Errorf("Merge: expected\n%v\n, actual \n%v\n", tt.m.ToString(), merged.ToString())
		}
	}
}

func TestAdd(t *testing.T) {
	for _, tt := range arithmeticTest {
		sum := matrix.Add(tt.m1, tt.m2)
		if !matrix.Equal(sum, tt.sum) {
			t.Errorf("Add: expected\n%v\n, actual \n%v\n", tt.sum.ToString(), sum.ToString())
		}
	}
}

func TestSub(t *testing.T) {
	for _, tt := range arithmeticTest {
		sub := matrix.Sub(tt.sum, tt.m1)
		if !matrix.Equal(sub, tt.m2) {
			t.Errorf("Sub: expected\n%v\n, actual \n%v\n", tt.m2.ToString(), sub.ToString())
		}
	}
}

func TestMul(t *testing.T) {
	for _, tt := range arithmeticTest {
		mul := matrix.Mul(tt.m1, tt.m2T)
		if !matrix.Equal(mul, tt.mul) {
			t.Errorf("Mul: expected\n%v\n, actual \n%v\n", tt.mul.ToString(), mul.ToString())
		}
	}
}

func TestSplitMultAddMerge(t *testing.T) {
	m1 := createMatrix(30, 50, func(i int) float64 { return float64(rand.Intn(10)) })
	m2 := createMatrix(50, 30, func(i int) float64 { return float64(rand.Intn(10)) })

	sm1 := m1.Split(20, 20)
	sm2 := m2.Split(20, 20)
	sizeOutput := len(sm1)
	subMul := make([][]*matrix.Matrix, sizeOutput)
	for i := range subMul {
		subMul[i] = make([]*matrix.Matrix, sizeOutput)
		for j := range subMul[i] {
			r, _ := sm1[i][0].Dims()
			_, c := sm2[0][j].Dims()
			subMul[i][j] = createMatrix(r, c, func(i int) float64 { return 0 })
			for k := range sm1[i] {
				p := matrix.Mul(sm1[i][k], sm2[k][j])
				subMul[i][j] = matrix.Add(subMul[i][j], p)
			}
		}
	}
	mul := matrix.Merge(subMul)

	res := matrix.Mul(m1, m2)
	if !matrix.Equal(res, mul) {
		t.Errorf("Mul: expected\n%v\n, actual \n%v\n", res.ToString(), mul.ToString())
	}
}
