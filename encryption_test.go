package pushbell

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryption(t *testing.T) {
	e, err := newEncryption()
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	uaPrivateKey, _ := ecdh.P256().GenerateKey(rand.Reader)
	uaPublicKey := uaPrivateKey.PublicKey()

	authSecret := make([]byte, 16)
	rand.Read(authSecret)

	auth := base64.RawURLEncoding.EncodeToString(authSecret)
	p256 := base64.RawURLEncoding.EncodeToString(uaPublicKey.Bytes())

	//salt := make([]byte, 16)
	//rand.Read(salt)

	//sharedSecret, err := e.privateKey.ECDH(uaPublicKey)

	e.encryptMessage(auth, p256, []byte("Test message"))

	//e.messageBody(uaPublicKey.Bytes(), authSecret, sharedSecret, salt, []byte("Some test message"))
}
