package main

// Go's encoding/asn1 package works directly with Go structs.
// No need for separate .asn schema files - the struct definition IS the schema.
// Go maps types automatically: string -> UTF8String, int -> INTEGER, []T -> SEQUENCE OF T

import (
	"encoding/asn1"
	"fmt"
	"os"
)

// Person defines our ASN.1 structure.
// This is equivalent to the ASN.1 definition:
//
//	Person ::= SEQUENCE {
//	    username        UTF8String,
//	    favouriteNumber INTEGER,
//	    interests       SEQUENCE OF UTF8String
//	}
type Person struct {
	Username        string
	FavouriteNumber int
	Interests       []string
}

func main() {
	// Create a Person instance with sample data
	person := Person{
		Username:        "Sebastian",
		FavouriteNumber: 7,
		Interests:       []string{"sailing", "reading", "coding"},
	}

	// Step 1: Encode struct to DER bytes using asn1.Marshal
	// This produces binary DER-encoded data following ASN.1 TLV format
	encoded, err := asn1.Marshal(person)
	if err != nil {
		fmt.Println("Error encoding:", err)
		return
	}

	fmt.Printf("Encoded %d bytes\n", len(encoded))
	fmt.Printf("Hex: %x\n", encoded)

	// Step 2: Save DER bytes to file
	// This creates a binary file that can be shared or inspected with tools like `xxd` or `dumpasn1`
	err = os.WriteFile("plain.der", encoded, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("Saved to plain.der")

	// Step 3: Read DER bytes from file
	// This simulates receiving DER data from an external source
	data, err := os.ReadFile("plain.der")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Printf("\nRead %d bytes from plain.der\n", len(data))

	// Step 4: Decode DER bytes back to struct using asn1.Unmarshal
	// This parses the TLV data and populates the struct fields
	var decoded Person
	_, err = asn1.Unmarshal(data, &decoded)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	fmt.Println("\nDecoded:")
	fmt.Printf("  Username: %s\n", decoded.Username)
	fmt.Printf("  FavouriteNumber: %d\n", decoded.FavouriteNumber)
	fmt.Printf("  Interests: %v\n", decoded.Interests)
}
