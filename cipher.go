package pushbell

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"errors"
)

const (
	maxPayloadBody = 4096
	maxPlainText   = 3993
)

var (
	ErrMessageTooLong = errors.New("message is too long")

	P256 = ecdh.P256()
)

func (api *API) cipherPlaintext(auth, p256dh string, message []byte) (*bytes.Buffer, error) {
	if len(message) > maxPlainText {
		return nil, ErrMessageTooLong
	}

	// Auth secret provided by the push service
	authSecret, err := parsBase64Key(auth)
	if err != nil {
		return nil, err
	}

	// User Agent public P256 key provided by the push service
	p256, err := parsBase64Key(p256dh)
	if err != nil {
		return nil, err
	}

	// User Agent public P256 key
	uaPublic, err := P256.NewPublicKey(p256)
	if err != nil {
		return nil, err
	}

	// Application Server private P256 key
	asPrivate, err := P256.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Application Server public P256 key
	asPublic := asPrivate.PublicKey()

	// Shared secret between the User Agent and the Application Server
	sharedKey, err := asPrivate.ECDH(uaPublic)
	if err != nil {
		return nil, err
	}

	// Salt
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	// Generate key_info according to RFC8291 3.4
	keyInfo := bytes.NewBuffer([]byte("WebPush: info\x00"))
	keyInfo.Write(uaPublic.Bytes())
	keyInfo.Write(asPublic.Bytes())

	// Generate PRK_key and IKM according to RFC8291 3.4:
	// First 3 arguments to hkdf.New() response for PRK
	//
	// PRK = HMAC-SHA-256(auth_secret, ecdh_secret)
	// IKM = HKDF-Expand(PRK, key_info, 32)
	ikm, err := hkdfExtractAndExpand(32, sharedKey, authSecret, keyInfo.Bytes())

	// Generate cek_info according to RFC8291 3.4:
	cekInfo := []byte("Content-Encoding: aes128gcm\x00")

	// Generate CEK according to RFC8291 3.4:
	// First 3 arguments to hkdf.New() response for PRK
	//
	// PRK = HMAC-SHA-256(salt, IKM)
	// CEK = HMAC-SHA-256(PRK, cek_info || 0x01)[0..15]
	cek, err := hkdfExtractAndExpand(16, ikm, salt, cekInfo)

	// Generate nonce_info according to RFC8291 3.4:
	nonceInfo := []byte("Content-Encoding: nonce\x00")

	// Generate NONCE according to RFC8291 3.4:
	// First 3 arguments to hkdf.New() response for PRK
	//
	// PRK = HMAC-SHA-256(salt, IKM)
	// NONCE = HMAC-SHA-256(PRK, nonce_info || 0x01)[0..11]
	nonce, err := hkdfExtractAndExpand(12, ikm, salt, nonceInfo)

	// Cipher block
	c, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}

	// GCM encryptor
	gcm, err := cipher.NewGCMWithNonceSize(c, 12)
	if err != nil {
		return nil, err
	}

	message = append(message, []byte("\x02")...)

	// Encrypt the payload
	ciphertext := gcm.Seal(nil, nonce, message, nil)

	// Record size greater than the sum of the lengths of the plaintext, the padding
	// delimiter (1 octet), any padding, and the authentication tag (16 octets).
	rs := make([]byte, 4)
	binary.LittleEndian.PutUint32(rs, maxPayloadBody)

	payloadBody := bytes.NewBuffer(salt)
	payloadBody.Write(rs)
	payloadBody.Write([]byte{65})
	payloadBody.Write(asPublic.Bytes())
	payloadBody.Write(ciphertext)

	return payloadBody, nil
}
