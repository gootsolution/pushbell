package pushbell

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type vapid struct {
	asSubject    string
	asPublicKey  string
	asPrivateKey *ecdsa.PrivateKey
}

var ErrSubjectNotValid = errors.New("subject VAPID should be either a \"mailto:\" (email) or a \"https:\" URI")

func newVAPID(asPrivateKey, asPublicKey, asSubject string) (*vapid, error) {
	if !strings.HasPrefix(asSubject, "mailto:") &&
		!strings.HasPrefix(asSubject, "https:") {
		return nil, ErrSubjectNotValid
	}

	publicKeyDecoded, err := parsBase64Key(asPublicKey)
	if err != nil {
		return nil, err
	}

	privateKeyDecoded, err := parsBase64Key(asPrivateKey)
	if err != nil {
		return nil, err
	}

	p256 := ecdh.P256()

	_, err = p256.NewPublicKey(publicKeyDecoded)
	if err != nil {
		return nil, err
	}

	privateECDH, err := p256.NewPrivateKey(privateKeyDecoded)
	if err != nil {
		return nil, err
	}

	privatePKCS8, err := x509.MarshalPKCS8PrivateKey(privateECDH)
	if err != nil {
		return nil, err
	}

	privateECDSA, err := x509.ParsePKCS8PrivateKey(privatePKCS8)
	if err != nil {
		return nil, err
	}

	return &vapid{
		asPublicKey:  asPublicKey,
		asPrivateKey: privateECDSA.(*ecdsa.PrivateKey),
	}, nil
}

const headerTemplate = `vapid t=%s, k=%s`

func (v *vapid) header(endpoint string) (string, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Subject:   v.asSubject,
		Audience:  jwt.ClaimStrings{fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
	})

	tokenSigned, err := token.SignedString(v.asPrivateKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(headerTemplate, tokenSigned, v.asPublicKey), nil
}
