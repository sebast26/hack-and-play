package main

import (
	"io"
	"os"
	"strings"
)

const CapitalLowerBound = 64
const CapitalUpperBound = 90
const SmallLowerBound = 96
const SmallUpperBound = 122
const Rot13 = 13

type rot13Reader struct {
	r io.Reader
}

func (reader rot13Reader) Read(bytes []byte) (int, error) {
	n, err := reader.r.Read(bytes)
	for i := range bytes {
		if isAlphaChar(bytes[i]) {
			bytes[i] = rot13(bytes[i])
		}
	}
	return n, err
}

func isAlphaChar(charVariable byte) bool {
	return (charVariable >= 'a' && charVariable <= 'z') || (charVariable >= 'A' && charVariable <= 'Z')
}

func rot13(b byte) byte {
	if b >= 'a' && b <= 'z' {
		newB := int(b) + Rot13
		if newB > SmallUpperBound {
			newB = SmallLowerBound + (newB % SmallUpperBound)
		}
		return byte(newB)
	}
	if b >= 'A' && b <= 'Z' {
		newB := int(b) + Rot13
		if newB > CapitalUpperBound {
			newB = CapitalLowerBound + (newB % CapitalUpperBound)
		}
		return byte(newB)
	}
	return b
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
