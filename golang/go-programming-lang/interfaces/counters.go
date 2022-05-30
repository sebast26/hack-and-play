package interfaces

import (
	"bufio"
	"bytes"
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c = ByteCounter(len(p))
	return len(p), nil
}

func testByteCounter() {
	var c ByteCounter
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // 12
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		*c++
	}
	return len(p), scanner.Err()
}

func testWordCounter() {
	var c WordCounter
	fmt.Fprintln(&c, "This is only a test")
	fmt.Println(c) // 5
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		*c++
	}
	return len(p), scanner.Err()
}

func testLineCounter() {
	var c LineCounter
	fmt.Fprintln(&c, "First line\nSecond line\nThird line")
	fmt.Println(c) // 3
}

func RunCounters() {
	testWordCounter()
	testLineCounter()
}
