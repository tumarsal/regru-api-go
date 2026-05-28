package regru

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"sort"
	"strings"
)

// SignatureAuth provides signature-based authentication for REG.RU API.
// This requires an SSL certificate with a private key uploaded to the API settings.
type SignatureAuth struct {
	Username   string
	PrivateKey *rsa.PrivateKey
}

// NewSignatureAuthFromFile creates signature authentication from a PEM private key file.
func NewSignatureAuthFromFile(username, privateKeyPath string) (*SignatureAuth, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	return NewSignatureAuth(username, data)
}

// NewSignatureAuth creates signature authentication from PEM-encoded private key data.
func NewSignatureAuth(username string, pemData []byte) (*SignatureAuth, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var key *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		var err error
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS1 private key: %w", err)
		}
	case "PRIVATE KEY":
		pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS8 private key: %w", err)
		}
		var ok bool
		key, ok = pk.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
	default:
		return nil, fmt.Errorf("unsupported PEM type: %s", block.Type)
	}

	return &SignatureAuth{
		Username:   username,
		PrivateKey: key,
	}, nil
}

// Sign creates a signature for the given parameters.
// The signature text is built by joining all parameter values (sorted by key, excluding "sig")
// with semicolons.
func (sa *SignatureAuth) Sign(params map[string]interface{}) (string, error) {
	sigText := makeTextForSig(params)

	hash := sha512.New()
	hash.Write([]byte(sigText))
	digest := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(nil, sa.PrivateKey, crypto.SHA512, digest)
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// makeTextForSig builds the signature text from parameters.
func makeTextForSig(params map[string]interface{}) string {
	parts := makeTextForSigRecursive(params)
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func makeTextForSigRecursive(v interface{}) []string {
	var results []string

	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			if k != "sig" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			results = append(results, makeTextForSigRecursive(val[k])...)
		}
	case []interface{}:
		for _, item := range val {
			results = append(results, makeTextForSigRecursive(item)...)
		}
	case []string:
		for _, item := range val {
			results = append(results, makeTextForSigRecursive(item)...)
		}
	case string:
		if val != "" {
			results = append(results, val)
		}
	case int:
		results = append(results, fmt.Sprintf("%d", val))
	case float64:
		results = append(results, fmt.Sprintf("%g", val))
	case bool:
		if val {
			results = append(results, "1")
		} else {
			results = append(results, "0")
		}
	}

	return results
}

// ClientWithSignature creates a new API client using signature authentication.
// Each request is signed with the RSA private key (zone/add_* and other methods).
func ClientWithSignature(username string, privateKey []byte, opts ...ClientOption) (*Client, error) {
	sa, err := NewSignatureAuth(username, privateKey)
	if err != nil {
		return nil, err
	}
	c := NewClient(username, "", opts...)
	c.sigAuth = sa
	return c, nil
}
