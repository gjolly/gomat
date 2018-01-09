package gomatcore_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/matei13/gomat/Daemon/gomatcore"
	"github.com/matei13/gomat/matrix"
)

func createMatrix(r, c int, fct func(int) float64) *matrix.Matrix {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = fct(i)
	}
	return matrix.New(r, c, data)
}

var splitMergeTest = []struct {
	m     *matrix.Matrix         // input matrix
	n     int                    // size of sub-matrices
	split []*gomatcore.SubMatrix // expected result
}{
	{
		createMatrix(5, 3, func(i int) float64 { return float64(i) }),
		2,
		[]*gomatcore.SubMatrix{
			&gomatcore.SubMatrix{
				Mat: matrix.New(2, 2, []float64{0, 1, 3, 4}),
				Row: 0,
				Col: 0,
			},
			&gomatcore.SubMatrix{
				Mat: matrix.New(2, 1, []float64{2, 5}),
				Row: 0,
				Col: 1,
			},
			&gomatcore.SubMatrix{
				Mat: matrix.New(2, 2, []float64{6, 7, 9, 10}),
				Row: 1,
				Col: 0,
			},
			&gomatcore.SubMatrix{
				Mat: matrix.New(2, 1, []float64{8, 11}),
				Row: 1,
				Col: 1,
			},
			&gomatcore.SubMatrix{
				Mat: matrix.New(1, 2, []float64{12, 13}),
				Row: 2,
				Col: 0,
			},
			&gomatcore.SubMatrix{
				Mat: matrix.New(1, 1, []float64{14}),
				Row: 2,
				Col: 1,
			},
		},
	},
}

func TestSplit(t *testing.T) {
	for _, tt := range splitMergeTest {
		split := gomatcore.Split(tt.m, tt.n)
		for i, result := range split {
			for j, expect := range tt.split {
				if (result.Row == expect.Row) && (result.Col == expect.Col) && !matrix.Equal(result.Mat, expect.Mat) {
					t.Errorf("Split(%d,%d): expected\n%v\n, actual \n%v\n", i, j, expect.Mat.ToString(), result.Mat.ToString())
				}
			}
		}
	}
}

func TestMerge(t *testing.T) {
	for _, tt := range splitMergeTest {
		r, c := tt.m.Dims()
		merged := gomatcore.Merge(tt.split, r, c, tt.n)
		if !matrix.Equal(merged, tt.m) {
			t.Errorf("Merge: expected\n%v\n, actual \n%v\n", tt.m.ToString(), merged.ToString())
		}
	}
}

func TestSplitMultAddMerge(t *testing.T) {
	m1 := createMatrix(30, 50, func(i int) float64 { return float64(rand.Intn(10)) })
	m2 := createMatrix(50, 30, func(i int) float64 { return float64(rand.Intn(10)) })
	blockSize := 20

	r, _ := m1.Dims()
	_, c := m2.Dims()

	nbRow := int(math.Ceil(float64(r) / float64(blockSize)))
	nbCol := int(math.Ceil(float64(c) / float64(blockSize)))

	sm1 := gomatcore.Split(m1, blockSize)
	sm2 := gomatcore.Split(m2, blockSize)

	subMul := make([][]*gomatcore.SubMatrix, nbRow*nbCol)
	for _, ssm1 := range sm1 {
		for _, ssm2 := range sm2 {
			if ssm1.Col == ssm2.Row {
				mul := matrix.Mul(ssm1.Mat, ssm2.Mat)
				subMul[ssm1.Row*nbCol+ssm2.Col] = append(subMul[ssm1.Row*nbCol+ssm2.Col], &gomatcore.SubMatrix{
					Mat: mul,
					Row: ssm1.Row,
					Col: ssm2.Col,
				})
			}
		}
	}

	matrices := make([]*gomatcore.SubMatrix, nbRow*nbCol)
	for i := range subMul {
		matrices[i] = subMul[i][0]
		for k := 1; k < len(subMul[i]); k++ {
			matrices[i].Mat = matrix.Add(matrices[i].Mat, subMul[i][k].Mat)
		}
	}
	mul := gomatcore.Merge(matrices, r, c, 20)

	res := matrix.Mul(m1, m2)
	if !matrix.Equal(res, mul) {
		t.Errorf("Mul: expected\n%v\n, actual \n%v\n", res.ToString(), mul.ToString())
	}
}
