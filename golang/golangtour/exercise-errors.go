package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %f", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	z := 1.0
	prevZ := 0.0
	iter := 0
	for math.Abs(z-prevZ) > 0.0000000001 {
		prevZ = z
		z -= (z*z - x) / (2 * z)
		iter++
	}
	fmt.Println("iterations:", iter)
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
