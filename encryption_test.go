package pushbell

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newEncryption(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "CreationOK",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newEncryption()
			if !tt.wantErr(t, err, "newEncryption()") {
				return
			}
		})
	}
}

func Test_encryption_rotate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "RotateOK",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, _ := newEncryption()
			tt.wantErr(t, e.rotate(), "rotate()")
		})
	}
}

func Test_encryption_encryptMessage(t *testing.T) {
	svc, _ := newEncryption()

	type args struct {
		auth    string
		p256dh  string
		message []byte
	}

	var tests = []struct {
		name    string
		service *encryption
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "EncryptionOK",
			service: svc,
			args: args{
				auth:    "rm_owGF0xliyVXsrZk1LzQ",
				p256dh:  "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y",
				message: []byte("{\"title\": \"My first message\"}"),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "EncryptionKeyError",
			service: svc,
			args: args{
				auth:    "rm_owGF0xliyVXssrZk1LzQ",
				p256dh:  "BKm5pKbGwkTxu7sdJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y",
				message: []byte("{\"title\": \"My first message\"}"),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.encryptMessage(tt.args.auth, tt.args.p256dh, tt.args.message)
			if !tt.wantErr(t, err, fmt.Sprintf("encryptMessage(%v, %v, %v)", tt.args.auth, tt.args.p256dh, tt.args.message)) {
				return
			}
		})
	}
}
