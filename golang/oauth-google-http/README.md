# Google OAuth 2.0 Service Account - Manual HTTP Implementation

A lightweight Go implementation of Google OAuth 2.0 authentication using Service Accounts, built with manual HTTP requests without any Google OAuth libraries.

## Overview

This project demonstrates how to authenticate with Google APIs using Service Accounts following the [Google Identity OAuth 2.0 Service Account guide](https://developers.google.com/identity/protocols/oauth2/service-account). It implements the complete OAuth 2.0 flow from scratch using only standard Go libraries and JWT signing.

## Key Features

- **No Google OAuth Library Dependencies**: Pure Go implementation using standard `net/http` and `crypto` packages
- **Service Account Authentication**: Server-to-server OAuth 2.0 flow without user interaction
- **JWT Construction**: Manual JWT header and claim construction with RS256 signing
- **PKCS8 Private Key Support**: Parse and use X.509 private keys from Google service account JSON files

## How It Works

1. **Load Service Account Key**: Reads the Google service account JSON file containing private key and credentials
2. **Create JWT**: Constructs a JWT token with proper claims (iss, aud, exp, iat, scope) and signs it with the private key
3. **Exchange JWT for Access Token**: Makes a POST request to Google's OAuth endpoint to exchange the JWT for an access token
4. **Use Access Token**: The returned access token can be used to authenticate requests to Google APIs

## Prerequisites

- Go 1.25 or higher
- Google Cloud Platform project with Service Account created
- Service Account JSON key file downloaded

## Installation

```bash
go get github.com/golang-jwt/jwt/v5
```

## Usage

Run the program with your service account JSON file:

```bash
go run main.go path/to/service-account.json
```

The program will:
1. Parse the service account credentials
2. Generate a JWT token signed with the private key
3. Request an access token from Google OAuth 2.0 endpoint
4. Display the access token and its metadata

### Example Output

```
JWT Token:
eyJhbGciOiJSUzI1NiIsImtpZCI6InRlc3Qta2V5LWlkIiwidHlwIjoiSldUIn0...

Access Token:
ya29.c.KqIBzwZ8...
Token Type: Bearer
Expires In: 3600 seconds
```

## Configuration

The default scope is set to `https://www.googleapis.com/auth/drive.readonly`. To change the scope, modify the `defaultScope` constant in `main.go`:

```go
const (
    oauthTokenURL = "https://oauth2.googleapis.com/token"
    defaultScope  = "https://www.googleapis.com/auth/your-desired-scope"
)
```

## Project Structure

```
.
├── main.go          # Main implementation
├── main_test.go     # Test suite
├── go.mod           # Go module definition
└── README.md        # This file
```

## Implementation Details

### Service Account Flow

This implementation uses Google Service Accounts, which enable server-to-server authentication without requiring user interaction. This is ideal for:
- Backend services
- Automated workflows
- Scheduled tasks
- Server applications

### JWT Token Structure

The JWT token includes:
- **Header**: Algorithm (RS256) and Key ID (kid)
- **Claims**: Issuer, Audience, Expiration, Issued At, and Scope
- **Signature**: RS256 signature using the service account's private key

### Manual HTTP Request

Instead of using Google's OAuth client libraries, this project makes a direct HTTP POST request to:
```
https://oauth2.googleapis.com/token
```

With form-encoded parameters:
- `grant_type`: `urn:ietf:params:oauth:grant-type:jwt-bearer`
- `assertion`: The signed JWT token

## Testing

Run the test suite:

```bash
go test -v
```

The tests cover:
- Private key parsing (PKCS8 format)
- Service account JSON loading
- JWT creation and signing
- Token request error handling

## References

- [Google OAuth 2.0 for Service Accounts](https://developers.google.com/identity/protocols/oauth2/service-account)
- [RFC 7519 - JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
- [RFC 7523 - JWT Bearer Token Grant Type](https://datatracker.ietf.org/doc/html/rfc7523)

## License

This project is provided as-is for educational and demonstration purposes.
