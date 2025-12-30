package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

const (
	charset               = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	pkceVerifierMinLength = 43
	pkceVerifierMaxLength = 128
)

// generatePKCECodeVerifier generates a cryptographically secure random string
// for use as a PKCE code verifier. The string length is randomly chosen between
// 43 and 128 characters and contains only alphanumeric characters (A-Z, a-z, 0-9)
// as required by RFC 7636.
func generatePKCECodeVerifier() (string, error) {
	lengthRange := pkceVerifierMaxLength - pkceVerifierMinLength + 1
	randomLength, err := rand.Int(rand.Reader, big.NewInt(int64(lengthRange)))
	if err != nil {
		return "", err
	}
	length := int(randomLength.Int64()) + pkceVerifierMinLength

	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// generatePKCECodeChallenge creates a code challenge from the given code verifier.
// It computes the BASE64-URL-encoded SHA-256 hash of the verifier as specified
// in RFC 7636 for the S256 challenge method.
func generatePKCECodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func main() {
	verifier, err := generatePKCECodeVerifier()
	if err != nil {
		fmt.Printf("Error generating code verifier: %v\n", err)
		return
	}

	challenge := generatePKCECodeChallenge(verifier)

	fmt.Printf("Code Verifier (length: %d): %s\n", len(verifier), verifier)
	fmt.Printf("Code Challenge: %s\n", challenge)
}
