package pushbell

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"sync"
)

type encryption struct {
	publicKey  []byte
	privateKey *ecdh.PrivateKey
	mu         sync.RWMutex
}

func newEncryption() (*encryption, error) {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey().Bytes()

	return &encryption{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (e *encryption) rotate() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	e.privateKey = privateKey
	e.publicKey = privateKey.PublicKey().Bytes()

	return nil
}

func (e *encryption) encryptMessage(auth, p256dh string, message []byte) (*bytes.Buffer, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	authSecretDecoded, err := parsBase64Key(auth)
	if err != nil {
		return nil, err
	}

	uaPublicKeyDecoded, err := parsBase64Key(p256dh)
	if err != nil {
		return nil, err
	}

	uaPublicKey, err := ecdh.P256().NewPublicKey(uaPublicKeyDecoded)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := e.privateKey.ECDH(uaPublicKey)
	if err != nil {
		return nil, err
	}

	salt := make([]byte, 16)

	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	recordSize := len(message) + 103

	buf := new(bytes.Buffer)
	buf.Grow(recordSize)

	e.messageHeader(buf, salt, uint32(recordSize))

	err = e.messageBody(buf, uaPublicKeyDecoded, authSecretDecoded, sharedSecret, salt, message)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (e *encryption) messageHeader(buf *bytes.Buffer, salt []byte, recordSize uint32) {
	// An 86-octet produced from the salt, record size of 4096, and application server public key
	header := buf.AvailableBuffer()[:86]

	// Writing salt
	copy(header, salt)

	// Writing record size as big endian uint32
	binary.BigEndian.PutUint32(header[16:21], recordSize)

	// Write key length (IDK why, in X9.62 it doesn't need)
	copy(header[20:], []byte{0x41})

	// Writing application server public key defined in X9.62
	copy(header[21:], e.publicKey)

	buf.Write(header)
}

func (e *encryption) messageBody(buf *bytes.Buffer, uaPublicKey, authSecret, sharedSecret, salt, message []byte) error {
	// Generate key_info according to RFC8291 3.4
	keyInfo := make([]byte, 144)
	copy(keyInfo[:14], "WebPush: info\x00")
	copy(keyInfo[14:79], uaPublicKey)
	copy(keyInfo[79:], e.publicKey)

	// Generate PRK_key and IKM according to RFC8291 3.4:
	// First 3 arguments to hkdf.New() response for PRK
	//
	// PRK = HMAC-SHA-256(auth_secret, ecdh_secret)
	// IKM = HKDF-Expand(PRK, key_info, 32)
	ikm, err := hkdfExtractAndExpand(32, sharedSecret, authSecret, keyInfo)

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
		return err
	}

	// GCM encryptor
	gcm, err := cipher.NewGCMWithNonceSize(c, 12)
	if err != nil {
		return err
	}

	mbuf := buf.AvailableBuffer()
	mbuf = append(mbuf, message...)
	mbuf = append(mbuf, byte(0x02))

	// Encrypt the payload
	sealed := gcm.Seal(nil, nonce, mbuf, nil)
	buf.Write(sealed)

	return nil
}
