package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("reports frequency of words in an input file\n\nusage: go run wordfreq.go <file name>")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("cannot open file: %v\n", err)
	}
	defer file.Close()

	wordfreq := make(map[string]int)

	input := bufio.NewScanner(file)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		wordfreq[input.Text()]++
	}

	for w, n := range wordfreq {
		fmt.Printf("%s\t%d\n", w, n)
	}
}
