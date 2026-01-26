package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Certificate struct {
	TBSCertificate     TBSCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type TBSCertificate struct {
	Version              int `asn1:"optional,explicit,tag:0,default:0"`
	SerialNumber         *big.Int
	Signature            pkix.AlgorithmIdentifier
	Issuer               pkix.RDNSequence
	Validity             Validity
	Subject              pkix.RDNSequence
	SubjectPublicKeyInfo asn1.RawValue
	Extensions           []pkix.Extension `asn1:"optional,explicit,tag:3"`
}

type Validity struct {
	NotBefore time.Time
	NotAfter  time.Time
}

func lookupOID(oid asn1.ObjectIdentifier) string {
	url := fmt.Sprintf("https://oidref.com/%s", oid.String())
	resp, err := http.Get(url)
	if err != nil {
		return oid.String()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return oid.String()
	}

	// Extract name from link pattern like "[15 keyUsage](/2.5.29.15)"
	re := regexp.MustCompile(`>\[\d+\s+([^\]]+)\]</a>`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		return string(matches[1])
	}

	// Fallback: extract name from title like "OID 2.5.29.15 keyUsage reference info"
	re2 := regexp.MustCompile(`<title>OID [\d.]+ ([^\s]+) reference`)
	matches2 := re2.FindSubmatch(body)
	if len(matches2) > 1 {
		return string(matches2[1])
	}

	return oid.String()
}

func main() {
	data, err := os.ReadFile("cert.der")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("Read %d bytes\n", len(data))

	var raw asn1.RawValue
	_, err = asn1.Unmarshal(data, &raw)
	if err != nil {
		fmt.Println("Not valid DER:", err)
		return
	}

	fmt.Println("Valid DER!")
	fmt.Printf("Tag: %d, Class: %d, IsCompound: %v\n", raw.Tag, raw.Class, raw.IsCompound)

	fmt.Println("\nChildren:")
	rest := raw.Bytes
	i := 0
	for len(rest) > 0 {
		var child asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &child)
		if err != nil {
			fmt.Println("Error parsing child:", err)
			break
		}
		fmt.Printf("  [%d] Tag: %d, Class: %d, IsCompound: %v, Len: %d\n",
			i, child.Tag, child.Class, child.IsCompound, len(child.Bytes))
		i++
	}

	var cert Certificate
	_, err = asn1.Unmarshal(data, &cert)
	if err != nil {
		fmt.Println("Error parsing certificate:", err)
		return
	}

	fmt.Println("\nParsed Certificate:")
	fmt.Printf("  Version: %d\n", cert.TBSCertificate.Version)
	fmt.Printf("  SerialNumber: %s\n", cert.TBSCertificate.SerialNumber)
	fmt.Printf("  Signature Algorithm: %v\n", cert.TBSCertificate.Signature.Algorithm)
	fmt.Printf("  Issuer: %s\n", cert.TBSCertificate.Issuer.String())
	fmt.Printf("  Subject: %s\n", cert.TBSCertificate.Subject.String())
	fmt.Printf("  NotBefore: %s\n", cert.TBSCertificate.Validity.NotBefore)
	fmt.Printf("  NotAfter: %s\n", cert.TBSCertificate.Validity.NotAfter)
	fmt.Printf("  SignatureValue len: %d bits\n", cert.SignatureValue.BitLength)

	fmt.Println("\nExtensions:")
	for _, ext := range cert.TBSCertificate.Extensions {
		name := lookupOID(ext.Id)
		fmt.Printf("  %s(%v)\n", name, ext.Id)
	}
}
