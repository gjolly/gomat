package matrix

type Operation int

const (
	OpAdd Operation = 1 + iota
	OpSub
	OpMul
)

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
