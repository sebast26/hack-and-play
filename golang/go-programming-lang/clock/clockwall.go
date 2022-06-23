package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type wallClock struct {
	location          string
	server            string
	currentTimeBuffer *bytes.Buffer
}

func main() {

	// parse args into wall clocks
	clocks := make([]wallClock, 0)
	for _, timeserver := range os.Args[1:] {
		tmp := strings.Split(timeserver, "=")
		clocks = append(clocks, wallClock{
			location:          tmp[0],
			server:            tmp[1],
			currentTimeBuffer: bytes.NewBuffer(make([]byte, 0)),
		})
	}

	for _, clock := range clocks {
		go netcat(clock)
	}

	var sb strings.Builder
	for {
		sb.Reset()
		for _, clock := range clocks {
			tmp := strings.Split(clock.currentTimeBuffer.String(), "\n")
			if len(tmp) > 1 {
				sb.WriteString(fmt.Sprintf("%s: %s ", clock.location, tmp[len(tmp)-2]))
			}
		}
		fmt.Printf("\r%s", sb.String())
		time.Sleep(1 * time.Second)
	}

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func netcat(clock wallClock) {
	conn, err := net.Dial("tcp", clock.server)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mustCopy(clock.currentTimeBuffer, conn)
}
