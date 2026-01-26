# ASN.1 and DER Encoding in Go

A learning project exploring ASN.1 (Abstract Syntax Notation One) encoding, DER format, and X.509 certificate parsing using Go's standard library.

## What is ASN.1?

ASN.1 is a standard for describing data structures in a platform-independent way. It's used extensively in:
- X.509 certificates (SSL/TLS)
- LDAP
- SNMP
- Cryptographic standards (PKCS)

## DER Encoding

DER (Distinguished Encoding Rules) is a binary encoding of ASN.1 data using **TLV (Tag-Length-Value)** format:

```
[Tag][Length][Value]
```

Each element starts with:
1. **Tag byte(s)** - identifies the data type
2. **Length byte(s)** - how many bytes in the value
3. **Value bytes** - the actual data

## Understanding Tags

### Class (2 bits)

| Value | Class | Meaning |
|-------|-------|---------|
| 0 | UNIVERSAL | Standard ASN.1 types (INTEGER, SEQUENCE, etc.) |
| 1 | APPLICATION | Application-specific types |
| 2 | CONTEXT-SPECIFIC | Context-dependent (used for optional fields, CHOICE) |
| 3 | PRIVATE | Organization-specific |

### IsCompound (1 bit)

- **false (primitive)** - contains actual data directly
- **true (constructed)** - contains nested TLV elements

### Common UNIVERSAL Tags

| Tag | Hex | Type |
|-----|-----|------|
| 1 | 0x01 | BOOLEAN |
| 2 | 0x02 | INTEGER |
| 3 | 0x03 | BIT STRING |
| 4 | 0x04 | OCTET STRING |
| 5 | 0x05 | NULL |
| 6 | 0x06 | OBJECT IDENTIFIER (OID) |
| 12 | 0x0C | UTF8String |
| 16 | 0x10 | SEQUENCE (0x30 with constructed bit) |
| 17 | 0x11 | SET (0x31 with constructed bit) |
| 19 | 0x13 | PrintableString |
| 22 | 0x16 | IA5String |
| 23 | 0x17 | UTCTime |
| 24 | 0x18 | GeneralizedTime |

## Go's encoding/asn1 Package

Go's `encoding/asn1` package works directly with Go structs - **no separate .asn schema files needed**.

### Type Mappings

| Go Type | ASN.1 Type |
|---------|------------|
| `int`, `int64` | INTEGER |
| `bool` | BOOLEAN |
| `string` | UTF8String |
| `[]byte` | OCTET STRING |
| `*big.Int` | INTEGER (arbitrary precision) |
| `time.Time` | UTCTime or GeneralizedTime |
| `asn1.ObjectIdentifier` | OBJECT IDENTIFIER |
| `asn1.BitString` | BIT STRING |
| `[]T` | SEQUENCE OF T |
| `struct` | SEQUENCE |

### Struct Tags

Control ASN.1 encoding with struct tags:

```go
type Example struct {
    Version    int    `asn1:"optional,explicit,tag:0,default:0"`
    Name       string `asn1:"utf8"`
    Items      []int  `asn1:"set"`
}
```

| Tag | Meaning |
|-----|---------|
| `optional` | Field may be absent |
| `explicit` | Wrapped in extra tag |
| `implicit` | Tag replaces default |
| `tag:N` | Use context-specific tag [N] |
| `default:V` | Default value if absent |
| `set` | Encode slice as SET (not SEQUENCE) |
| `utf8` | Force UTF8String encoding |

### asn1.RawValue

For inspecting unknown structures:

```go
type RawValue struct {
    Class, Tag int
    IsCompound bool
    Bytes      []byte  // Content only
    FullBytes  []byte  // Tag + Length + Content
}
```

## Object Identifiers (OIDs)

OIDs are hierarchical identifiers like `2.5.4.3` (CommonName).

### Common X.500 Attribute OIDs (2.5.4.x)

| OID | Name |
|-----|------|
| 2.5.4.3 | CN (CommonName) |
| 2.5.4.6 | C (Country) |
| 2.5.4.7 | L (Locality) |
| 2.5.4.8 | ST (State) |
| 2.5.4.10 | O (Organization) |
| 2.5.4.11 | OU (OrganizationalUnit) |

### Common Extension OIDs (2.5.29.x)

| OID | Name |
|-----|------|
| 2.5.29.14 | subjectKeyIdentifier |
| 2.5.29.15 | keyUsage |
| 2.5.29.17 | subjectAltName |
| 2.5.29.19 | basicConstraints |
| 2.5.29.31 | cRLDistributionPoints |
| 2.5.29.32 | certificatePolicies |
| 2.5.29.35 | authorityKeyIdentifier |
| 2.5.29.37 | extKeyUsage |

### Looking Up OIDs

- http://oid-info.com
- https://oidref.com/2.5.29.17

## X.509 Certificate Structure

```
Certificate ::= SEQUENCE {
    tbsCertificate       TBSCertificate,
    signatureAlgorithm   AlgorithmIdentifier,
    signatureValue       BIT STRING
}

TBSCertificate ::= SEQUENCE {
    version         [0]  EXPLICIT Version DEFAULT v1,
    serialNumber         CertificateSerialNumber,
    signature            AlgorithmIdentifier,
    issuer               Name,
    validity             Validity,
    subject              Name,
    subjectPublicKeyInfo SubjectPublicKeyInfo,
    extensions      [3]  EXPLICIT Extensions OPTIONAL
}
```

### Extension Critical Flag

- **Critical: true** - Validator MUST understand this extension or reject the certificate
- **Critical: false** - Validator MAY ignore if not understood

## Go Packages for Certificates

### encoding/asn1
Low-level ASN.1 marshaling/unmarshaling.

### crypto/x509/pkix
Building blocks for X.509:
- `RDNSequence`, `RelativeDistinguishedNameSET`
- `AlgorithmIdentifier`
- `Extension`
- `Name` (with parsed CN, O, C fields)

### crypto/x509
Full certificate parsing:
```go
cert, err := x509.ParseCertificate(derBytes)
// Access cert.Subject.CommonName, cert.DNSNames, cert.NotAfter, etc.
```

---

## Programs in This Repository

### plain.go - ASN.1 Encoding Basics

Demonstrates encoding a simple Go struct to DER format and decoding it back.

**What it does:**
1. Defines a `Person` struct with username, favourite number, and interests
2. Encodes it to DER bytes using `asn1.Marshal`
3. Saves the DER bytes to `plain.der`
4. Reads the file back
5. Decodes using `asn1.Unmarshal`

**Run:**
```bash
go run plain.go
```

**Inspect the output:**
```bash
xxd plain.der
```

### cert.go - X.509 Certificate Parser

Parses a DER-encoded X.509 certificate and displays its structure.

**What it does:**
1. Reads `cert.der` file
2. Validates it's proper DER format
3. Shows raw TLV structure (tag, class, compound, length)
4. Parses into custom Certificate/TBSCertificate structs
5. Displays parsed fields (version, serial, issuer, subject, validity)
6. Lists extensions with OID names (fetched from oidref.com)

**Run:**
```bash
go run cert.go
```

**Requirements:**
- Place a DER-encoded certificate as `cert.der` in the same directory
- Internet connection for OID lookups (optional, falls back to OID numbers)

**Convert PEM to DER:**
```bash
openssl x509 -in cert.pem -outform DER -out cert.der
```

---

## Useful Commands

```bash
# View DER file as hex
xxd file.der

# Parse DER with OpenSSL
openssl asn1parse -inform DER -in file.der

# View certificate details
openssl x509 -inform DER -in cert.der -text -noout

# Convert PEM to DER
openssl x509 -in cert.pem -outform DER -out cert.der

# Convert DER to PEM
openssl x509 -inform DER -in cert.der -out cert.pem
```

## References

- [ITU-T X.680](https://www.itu.int/rec/T-REC-X.680) - ASN.1 Specification
- [ITU-T X.690](https://www.itu.int/rec/T-REC-X.690) - DER/BER Encoding Rules
- [RFC 5280](https://tools.ietf.org/html/rfc5280) - X.509 Certificate Profile
- [Go encoding/asn1](https://pkg.go.dev/encoding/asn1)
- [Go crypto/x509](https://pkg.go.dev/crypto/x509)
- [OID Repository](http://oid-info.com)
