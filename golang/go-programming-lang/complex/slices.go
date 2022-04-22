package complex

import (
	"bytes"
	"fmt"
	"reflect"
)

func slice() {
	months := [...]string{1: "January", 2: "February", 3: "March", 4: "April", 5: "May", 6: "June", 7: "July",
		8: "August", 9: "September", 10: "October", 11: "November", 12: "December"}
	Q2 := months[4:7]
	summer := months[6:9]
	fmt.Println(Q2)
	fmt.Println(summer)

	// A slice has three components: a pointer, a length, and a capacity.
	// 1. The pointer points to the first element of the array that is reachable through the slice, which is not necessarily the array’s first element.
	// 2. The length is the number of slice elements; it can’t exceed the capacity,
	// 3. capacity is usually the number of elements between the start of the slice and the end of the underlying array.

	fmt.Printf("%v len=%d cap=%d\n", months[:], len(months[:]), cap(months[:]))
	fmt.Printf("%v len=%d cap=%d\n", Q2, len(Q2), cap(Q2))
}

func slicesNotComparable() {
	a := make([]byte, 10)
	b := make([]byte, 10)

	// fmt.Println(a == b) // invalid operation
	fmt.Println(bytes.Equal(a, b))
	a = append(a, 1)
	fmt.Println(bytes.Equal(a, b))
}

func sliceCompare() {
	a := []int{9, 10, 11}
	b := []int{9, 10, 11}

	fmt.Println(reflect.DeepEqual(a, b))
	a = append(a, 12)
	fmt.Println(reflect.DeepEqual(a, b))
}

func makeSlice() {
	_ = make([]string, 10)     // len; the slice is a view for whole array
	_ = make([]string, 10, 20) // len, cap
}

// The built-in append function may use a more sophisticated growth strategy than appendInt’s simplistic one
func appendInt(x []int, y ...int) []int {
	var z []int
	zlen := len(x) + len(y)
	if zlen <= cap(x) {
		// There is room to grow.  Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.  Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // a built-in function; see text
	}
	copy(z[len(x):], y)
	return z
}

func expandingSlice() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("%d  cap=%d\t%v\n", i, cap(y), y)
		x = y
	}

	//0  cap=1        [0]
	//1  cap=2        [0 1]
	//2  cap=4        [0 1 2]
	//3  cap=4        [0 1 2 3]
	//4  cap=8        [0 1 2 3 4]
	//5  cap=8        [0 1 2 3 4 5]
	//6  cap=8        [0 1 2 3 4 5 6]
	//7  cap=8        [0 1 2 3 4 5 6 7]
	//8  cap=16       [0 1 2 3 4 5 6 7 8]
	//9  cap=16       [0 1 2 3 4 5 6 7 8 9]
}

func multiAppend() {
	var x []int
	x = append(x, 1)
	x = append(x, 2, 3)
	x = append(x, 4, 5, 6)
	x = append(x, x...) // append the slice x
	fmt.Println(x)      // "[1 2 3 4 5 6 1 2 3 4 5 6]"
}

func RunSlice() {
	//slice()
	//slicesNotComparable()
	//sliceCompare()
	//makeSlice()
	//expandingSlice()
	multiAppend()
}
