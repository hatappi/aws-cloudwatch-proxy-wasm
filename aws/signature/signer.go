package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	signatureV4Algorithm   = "AWS4-HMAC-SHA256"
	signatureV4ContentType = "application/x-amz-json-1.0"
)

// SignatureV4Config is config for SignatureV4
type SignatureV4Config struct {
	Region      string
	Service     string
	Method      string
	Host        string
	Path        string
	Querystring string
	Body        []byte
}

// Signer defines the methods required for signature
type Signer interface {
	SetSignatureV4Header(header http.Header, config SignatureV4Config) error
}

// New initializes signer that meets Signer interface
func New(accessKeyID, secretAccessKey string) Signer {
	return &signer{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
}

type signer struct {
	accessKeyID     string
	secretAccessKey string
}

// SetSignatureV4Header sets SignatureV4 header to header specified by the argument
func (s *signer) SetSignatureV4Header(header http.Header, config SignatureV4Config) error {
	if header == nil {
		return errors.New("header is nil")
	}

	now := time.Now().UTC()
	amzdate := now.Format("20060102T150405Z")
	datestamp := now.Format("20060102")

	//
	// create a canonical request
	//
	canonicalHeaderMap := map[string]string{
		"content-type": signatureV4ContentType,
		"host":         config.Host,
		"x-amz-date":   amzdate,
	}
	for k := range header {
		canonicalHeaderMap[strings.ToLower(k)] = header.Get(k)
	}

	var canonicalHeaderKeys []string
	for k := range canonicalHeaderMap {
		canonicalHeaderKeys = append(canonicalHeaderKeys, k)
	}
	sort.Strings(canonicalHeaderKeys)
	signedHeader := strings.Join(canonicalHeaderKeys, ";")

	var canonicalHeader string
	for _, k := range canonicalHeaderKeys {
		canonicalHeader += fmt.Sprintf("%s:%s\n", k, canonicalHeaderMap[k])
	}

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", config.Method, config.Path, config.Querystring, canonicalHeader, signedHeader, toHash(config.Body))

	//
	// create the string to sign
	//
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", datestamp, config.Region, config.Service)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", signatureV4Algorithm, amzdate, credentialScope, toHash([]byte(canonicalRequest)))

	//
	// calculate the signature
	//
	signingKey := generateSignatureKey(s.secretAccessKey, datestamp, config.Region, config.Service)
	signature := hex.EncodeToString(sign(signingKey, stringToSign))

	//
	// add signing information to the request
	//
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", signatureV4Algorithm, s.accessKeyID, credentialScope, signedHeader, signature)

	header.Set("Content-Type", signatureV4ContentType)
	header.Set("x-amz-date", amzdate)
	header.Set("Authorization", authorizationHeader)

	return nil
}

func toHash(data []byte) []byte {
	d := sha256.Sum256(data)
	return []byte(hex.EncodeToString(d[:]))
}

func sign(key []byte, msg string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	return h.Sum(nil)
}

func generateSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := sign([]byte("AWS4"+key), dateStamp)
	kRegion := sign(kDate, regionName)
	kService := sign(kRegion, serviceName)
	kSigning := sign(kService, "aws4_request")
	return kSigning
}
