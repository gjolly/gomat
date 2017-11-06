package matrix

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func createMatrix(r, c int) *Matrix {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		data[i] = float64(i)
	}
	return New(r, c, data)
}

func format(m *Matrix) fmt.Formatter {
	return mat.Formatted(m, mat.Squeeze())
}

var splitMergeTest = []struct {
	m     *Matrix     // input matrix
	i     int         // nb of sub-columns
	j     int         // nb of sub-rows
	split [][]*Matrix // expected result
}{
	{
		createMatrix(5, 5), 2, 2, [][]*Matrix{
			[]*Matrix{
				New(3, 3, []float64{0, 1, 2, 5, 6, 7, 10, 11, 12}),
				New(3, 2, []float64{3, 4, 8, 9, 13, 14}),
			},
			[]*Matrix{
				New(2, 3, []float64{15, 16, 17, 20, 21, 22}),
				New(2, 2, []float64{18, 19, 23, 24}),
			},
		},
	},
}

func TestSplit(t *testing.T) {
	for _, tt := range splitMergeTest {
		split := (tt.m).Split(tt.i, tt.j)
		for i := range split {
			for j := range split[i] {
				if !mat.Equal(split[i][j], tt.split[i][j]) {
					t.Errorf("Split(%d,%d): expected\n%v\n, actual \n%v\n", i, j, format(tt.split[i][j]), format(split[i][j]))
				}
			}
		}
	}
}

func TestMerge(t *testing.T) {
	for _, tt := range splitMergeTest {
		merged := Merge(tt.split)
		if !mat.Equal(merged, tt.m) {
			t.Errorf("Merge: expected\n%v\n, actual \n%v\n", format(tt.m), format(merged))
		}
	}
}
