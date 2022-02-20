package main

import (
	"gdoc-writer/google"
	"gdoc-writer/stdin"
	"log"
)

func main() {
	b, err := stdin.ReadStdin()
	if err != nil {
		log.Fatalf("Unable to read data from stdin: %v", err)
	}

	docLocation := google.CreateDocument(string(b))
	log.Printf("Successfully created Google Document under: %s", docLocation)
}
