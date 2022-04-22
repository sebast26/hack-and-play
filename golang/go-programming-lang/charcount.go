package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func inc(counts map[string]map[rune]int, r rune, s string) {
	if _, ok := counts[s]; !ok {
		counts[s] = make(map[rune]int)
	}
	counts[s][r]++
}

func main() {
	counts := make(map[string]map[rune]int) // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int         // count of lengths of UTF-8 encodings
	invalid := 0                            // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		switch {
		case unicode.IsLetter(r):
			inc(counts, r, "letter")
		case unicode.IsSpace(r):
			inc(counts, r, "space")
		case unicode.IsNumber(r):
			inc(counts, r, "number")
		}
		utflen[n]++
	}
	fmt.Printf("cat\trune\tcount\n")
	for cat, m := range counts {
		fmt.Println(cat)
		for c, n := range m {
			fmt.Printf("\t%q\t%d\n", c, n)
		}
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
