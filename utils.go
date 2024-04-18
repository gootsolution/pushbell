package pushbell

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/hkdf"
)

var ErrKeyNotValid = errors.New("key is not a valid ECDSA private key")

func parseApplicationKeys(publicKey, privateKey *string) (*string, *ecdsa.PrivateKey, error) {
	_, err := parsBase64Key(*publicKey)
	if err != nil {
		return nil, nil, err
	}

	privateDecoded, err := parsBase64Key(*privateKey)
	if err != nil {
		return nil, nil, err
	}

	privateECDH, err := ecdh.P256().NewPrivateKey(privateDecoded)
	if err != nil {
		return nil, nil, err
	}

	sourceEncoded, err := x509.MarshalPKCS8PrivateKey(privateECDH)
	if err != nil {
		return nil, nil, err
	}

	ecdsaPrivate, err := x509.ParsePKCS8PrivateKey(sourceEncoded)
	if err != nil {
		return nil, nil, err
	}

	privateVAPID, ok := ecdsaPrivate.(*ecdsa.PrivateKey)
	if !ok {
		return nil, nil, ErrKeyNotValid
	}

	return publicKey, privateVAPID, nil
}

func parsBase64Key(keyString string) ([]byte, error) {
	if strings.Contains(keyString, "=") {
		return base64.URLEncoding.DecodeString(keyString)
	} else {
		return base64.RawURLEncoding.DecodeString(keyString)
	}
}

func hkdfExtractAndExpand(length int, secret, salt, info []byte) ([]byte, error) {
	buf := make([]byte, length)

	reader := hkdf.New(sha256.New, secret, salt, info)

	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	return buf, nil
}
