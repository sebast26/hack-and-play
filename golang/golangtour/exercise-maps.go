package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	ret := make(map[string]int)
	fields := strings.Fields(s)
	for _, f := range fields {
		ret[f] += 1
	}
	return ret
}

func main() {
	wc.Test(WordCount)
}

