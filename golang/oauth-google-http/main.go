package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	oauthTokenURL = "https://oauth2.googleapis.com/token"
	defaultScope  = "https://www.googleapis.com/auth/drive.readonly"
)

type ServiceAccountKey struct {
	PrivateKey   string `json:"private_key"`
	PrivateKeyID string `json:"private_key_id"`
	ClientEmail  string `json:"client_email"`
	ProjectID    string `json:"project_id"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

func CreateJWT(privateKey *rsa.PrivateKey, keyID string, claims CustomClaims) (string, error) {
	jwt.MarshalSingleStringAsArray = false
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ParsePrivateKeyFromPEM(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("key is not RSA private key")
	}

	return nil, fmt.Errorf("failed to parse private key")
}

func LoadServiceAccountKey(filePath string) (*ServiceAccountKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	var key ServiceAccountKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &key, nil
}

func RequestAccessToken(jwtToken string) (*TokenResponse, error) {
	formData := url.Values{}
	formData.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	formData.Set("assertion", jwtToken)

	resp, err := http.PostForm(oauthTokenURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &tokenResp, nil
}

func GetAccessToken(serviceAccount *ServiceAccountKey) (*TokenResponse, error) {
	privateKey, err := ParsePrivateKeyFromPEM([]byte(serviceAccount.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    serviceAccount.ClientEmail,
			Audience:  jwt.ClaimStrings{oauthTokenURL},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Scope: defaultScope,
	}

	jwtToken, err := CreateJWT(privateKey, serviceAccount.PrivateKeyID, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT: %w", err)
	}

	return RequestAccessToken(jwtToken)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-service-account-json>")
	}

	keyFilePath := os.Args[1]

	serviceAccount, err := LoadServiceAccountKey(keyFilePath)
	if err != nil {
		log.Fatalf("Failed to load service account key: %v", err)
	}

	accessTokenResp, err := GetAccessToken(serviceAccount)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	fmt.Println("\nAccess Token:")
	fmt.Println(accessTokenResp.AccessToken)
	fmt.Printf("Token Type: %s\n", accessTokenResp.TokenType)
	fmt.Printf("Expires In: %d seconds\n", accessTokenResp.ExpiresIn)
}
