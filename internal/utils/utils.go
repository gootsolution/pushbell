package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/hkdf"
)

func ParseBase64Key(keyString string) ([]byte, error) {
	padded := strings.Contains(keyString, "=")

	var enc *base64.Encoding

	switch {
	case strings.ContainsAny(keyString, "+/"):
		enc = base64.RawStdEncoding
		if padded {
			enc = base64.StdEncoding
		}
	case strings.ContainsAny(keyString, "-_"):
		enc = base64.RawURLEncoding
		if padded {
			enc = base64.URLEncoding
		}
	default:
		return nil, errors.New("invalid base64 string")
	}

	return enc.DecodeString(keyString)
}

func HkdfExtractAndExpand(length int, secret, salt, info []byte) ([]byte, error) {
	buf := make([]byte, length)

	reader := hkdf.New(sha256.New, secret, salt, info)

	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	return buf, nil
}
