package main

import (
	"flag"
	complex2 "go-programming-lang/complex"
	"log"
	"os"
)

var sha384 = flag.Bool("sha384", false, "use SHA384 hash function")
var sha512 = flag.Bool("sha512", false, "use SHA512 hash function")

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: go run complex.go <input>")
	}

	input := os.Args[1]

	flag.Parse()

	complex2.Run(input, sha384, sha512)

	// ------

	//complex2.RunSlice()
	complex2.RunMaps()

}
