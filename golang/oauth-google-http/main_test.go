package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testPrivateKeyPKCS8 = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC7pMBVTXdrhJGv
fxfSXiGcOUeLhSfZg6ilLpnKUn0j89J7wCPCKbDt1BAlqk0H9P2hqHe78PxWOBFa
DprRuOjJt0FEpxYJw2dgpJGMOw6HWD9rTuGhp6Py0NtgxIROXKlnItZT4BrUzU9k
pyz8QIp2vIYhvu9OwSUCAruSzxr/iz7IEclEo/n7MUZ12NcvkR6d7tuPyNnqXEkE
ivauEeTNGePs0mG+qAa8MF1/YTCYFBjv9AppZJ8b3W4POnfAm6nWZ4fbQwOStfpy
hE1vbPy8nzM7GqmipGLzTJGf/uvD4F6TI9R+YmvkpWdDyk5Zp5TkNlrYAWSYG3ob
hofx3lGZAgMBAAECggEAAv/4ijSVmOyy5o8d2qFiIK0CwY6Eq1PRRajH0SX99a3v
esaA29kOSNMX2f5NwneOeX/nwCOzaQ526o1quLkf0Z/CJZ6VOd07oR4obIlCK3RD
+Ut2zH6RXq+rnjGuGtGc7LbwHcGqtCW/I/EPRU4j2I+dhlsW2Ne6J3zTMl7ow9aX
Qu8qNLdvckdwDQuc2vg9N6wrI95U+GFA6snuTz1lmBj2wD9Kf7/oadzwrRIpF8FA
pa04cHdjVwgrcskLAvrY8hXv1OrrUvjh7WDxQYZJ/6KQfwUzpTE1P4pwbdNwwA50
t4nmrvz5t0DWsuJK4EJjcfn96vldOF8ZsyaIJZl/PQKBgQD0XDc83u6XfntF4a0R
zdohyIIhWYeX74kzveu6V9rWMmr75PDOyvS2frGQgtTLMKwYd9hhN1mFh5NVTA13
MMs3BvzJTEotU0ie95x7wRo/HjcCobtEBLxJbH6ROoyCG83h2IZbmc12v/kTjZCc
sH5xq2bkdnPP3Mi7uvklGO+ydQKBgQDElOtvKtnPxhmrrVu+4wFHvk1QtkrTfU4I
QTkE9CBAof/BrhQXP5avxPPjR3n4VHJ4k4F1Welt5Mt51a+l2Pch5CclNsR9lNUK
TDVsSnDpAabj6zypQwzyjkWrytDj0ub9fmtUZGPrQMIVmMkN0qH4rgtpMgW+0lWY
3y6yClA2FQKBgQDHIAJPUydAQmx39MAQ2xufHgKU5J0u/Z5jkKdwUTleBPXTUKu/
B/pEYJeXuFVpQr9qrPahufX0bxdrka8Darqem84Qx19G8jAigop8+k3GE3wmY97z
PZylhIjbFcf8GmY6XTSXnoyk8SkqetjWJqjTPiZ5k+EqCLN275UpWtmxJQKBgQCJ
zKpMvwj38LHF+ju5z37yP6AK1/4Tyl6MwsupgUeqhuS1a8WBln2WsMUVqOEDBo9H
nn69OE29TwijWvtJkpGKHFikaYPLVC0SjuFUC8qiSxol4pwfOnU1Ai4sgntzeD+b
qqap/cVc/4QNh5BINEJ+A+LT1tDYNrDx/GDLRRd9EQKBgQDN8H7W6Iu+rog82BkY
rjFUqUc6XRpi2bzc1qwmOmzM4Tqozc0rb3eufs2hhP+Zu6ocVOuAmRQvVBXTEe38
E8ikwrFMH0B2t/kbol6HhN1+V/2G9igLZ7jtTCWfyWpbesucLFctBkT6w5hK2bh/
YudmuupazmQPUXhf25HS8j+FTQ==
-----END PRIVATE KEY-----`

func TestParsePrivateKeyFromPEM(t *testing.T) {
	t.Run("valid PKCS8 key", func(t *testing.T) {
		key, err := ParsePrivateKeyFromPEM([]byte(testPrivateKeyPKCS8))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if key == nil {
			t.Fatal("expected key, got nil")
		}
	})

	t.Run("invalid PEM data", func(t *testing.T) {
		_, err := ParsePrivateKeyFromPEM([]byte("not a valid pem"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty data", func(t *testing.T) {
		_, err := ParsePrivateKeyFromPEM([]byte(""))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestLoadServiceAccountKey(t *testing.T) {
	t.Run("valid JSON file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "service-account-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		testData := ServiceAccountKey{
			PrivateKey:   testPrivateKeyPKCS8,
			PrivateKeyID: "test-key-id",
			ClientEmail:  "test@example.com",
			ProjectID:    "test-project",
		}
		data, _ := json.Marshal(testData)
		tmpFile.Write(data)
		tmpFile.Close()

		key, err := LoadServiceAccountKey(tmpFile.Name())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if key.ClientEmail != "test@example.com" {
			t.Errorf("expected client_email test@example.com, got %s", key.ClientEmail)
		}
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.Write([]byte("invalid json"))
		tmpFile.Close()

		_, err = LoadServiceAccountKey(tmpFile.Name())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadServiceAccountKey("/non/existent/file.json")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCreateJWT(t *testing.T) {
	privateKey, err := ParsePrivateKeyFromPEM([]byte(testPrivateKeyPKCS8))
	if err != nil {
		t.Fatalf("failed to parse test key: %v", err)
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test@example.com",
			Audience:  jwt.ClaimStrings{oauthTokenURL},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Scope: defaultScope,
	}

	token, err := CreateJWT(privateKey, "test-key-id", claims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if kid, ok := parsedToken.Header["kid"].(string); !ok || kid != "test-key-id" {
		t.Errorf("expected kid test-key-id, got %v", parsedToken.Header["kid"])
	}
}

func TestRequestAccessToken(t *testing.T) {
	t.Run("empty JWT", func(t *testing.T) {
		_, err := RequestAccessToken("")
		if err == nil {
			t.Fatal("expected error with empty JWT, got nil")
		}
	})
}

func TestGetAccessToken(t *testing.T) {
	t.Run("with invalid private key", func(t *testing.T) {
		serviceAccount := &ServiceAccountKey{
			PrivateKey:   "invalid key",
			PrivateKeyID: "test-key-id",
			ClientEmail:  "test@example.com",
			ProjectID:    "test-project",
		}

		_, err := GetAccessToken(serviceAccount)
		if err == nil {
			t.Fatal("expected error with invalid private key, got nil")
		}
	})
}
