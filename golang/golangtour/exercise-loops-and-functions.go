package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	prevZ := 0.0
	iter := 0
	for math.Abs(z - prevZ) > 0.0000000001 {
		prevZ = z
		z -= (z*z - x) / (2*z)
		iter++
	}
	fmt.Println("iterations:", iter)
	return z
}

func main() {
	fmt.Println(Sqrt(3))
}
