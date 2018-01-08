package main

import (
	"fmt"
	"time"

	"github.com/matei13/gomat/gomat"
)

func main() {
	m1 := gomat.New(3, 2, []float64{0, 1, 2, 3, 4, 5})
	m2 := gomat.New(3, 2, []float64{1, 2, 4, 6, 8, 10})

	start := time.Now()
	gomat.Mult(m1, m2)
	elapsed := time.Since(start)
	fmt.Printf("Multiplication computed in %s.", elapsed.String())
}
