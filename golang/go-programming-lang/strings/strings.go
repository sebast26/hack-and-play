package strings

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"
)

// A string is an immutable sequence of bytes. Text strings are conventionally interpreted as UTF-8-encoded sequences of
// Unicode code points (runes).

// UTF-8 is a variable-length encoding of Unicode code points as bytes (runes). It uses between 1 and 4 bytes to represent
// each rune, but only 1 byte for ASCII characters.

func unicodeCount() {
	s := "Hello, 世界"
	fmt.Println(len(s))                    // 13
	fmt.Println(utf8.RuneCountInString(s)) // 9
}

func useRangeToIterateOverString() {
	for i, r := range "Hello, 世界" {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}

	//0       'H'     72
	//1       'e'     101
	//2       'l'     108
	//3       'l'     108
	//4       'o'     111
	//5       ','     44
	//6       ' '     32
	//7       '世'    19990
	//10      '界'    30028
}

// replacementCharacter is used whenever UTF-8 Decoder consumes an unexpected input byte
func replacementCharacter() {
	fmt.Printf("%s\n", string('\uFFFD')) // this is replacement character
	fmt.Printf("%s\n", string(1234567))  // invalid rune
}

func commaRec(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return commaRec(s[:n-3]) + "," + s[n-3:]
}

func comma(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		if (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func strToInt(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("strToInt error: %v", err)
	}
	return x
}

func Run() {
	unicodeCount()
	useRangeToIterateOverString()
	replacementCharacter()

	fmt.Println(commaRec("12345"))
	fmt.Println(comma("12345"))
	fmt.Println(commaRec("12345678901234567890"))
	fmt.Println(comma("12345678901234567890"))

	fmt.Println(strToInt("123456789"))
}
