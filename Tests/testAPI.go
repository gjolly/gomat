package main

import (
	"github.com/matei13/gomat/gomat"
	"github.com/matei13/gomat/matrix"
	"fmt"
)

func main() {
	m1 := matrix.New(2, 2, []float64{2, 2, 2, 2})
	m2 := matrix.New(2, 2, []float64{1, 1, 1, 1})
	r, err := gomat.Add(m1, m2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sum:")
	fmt.Println(r)

	r, err = gomat.Sub(m1, m2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sub:")
	fmt.Println(r)

	r, err = gomat.Mult(m1, m2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mult:")
	fmt.Println(r)
}
