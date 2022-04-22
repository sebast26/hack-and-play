package complex

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func arrays() {
	var a [5]int
	var b [6]int
	c := [...]int{1, 2, 3}

	fmt.Printf("%T\t%T\t%T", a, b, c)
}

func hashes(s string, useSHA384 *bool, useSHA512 *bool) {
	switch {
	case *useSHA384:
		fmt.Printf("%x\n", sha512.Sum384([]byte(s)))
	case *useSHA512:
		fmt.Printf("%x\n", sha512.Sum512([]byte(s)))
	default:
		fmt.Printf("%x\n", sha256.Sum256([]byte(s)))
	}
}

func Run(s string, sha384 *bool, sha512 *bool) {
	hashes(s, sha384, sha512)
}
