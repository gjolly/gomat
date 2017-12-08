package matrix_test

import (
	"testing"

	"../matrix"
)

func createMatrix(r, c int) *matrix.Matrix {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = float64(i)
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
		createMatrix(5, 3), 3, 2, [][]*matrix.Matrix{
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
