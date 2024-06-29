package pushbell

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"

	"golang.org/x/crypto/hkdf"
)

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
