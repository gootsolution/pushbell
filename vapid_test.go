package pushbell

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newVAPID(t *testing.T) {
	type args struct {
		asPublicKey  string
		asPrivateKey string
		asSubject    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "NewVapidOK",
			args: args{
				asPrivateKey: "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
				asPublicKey:  "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
				asSubject:    "mailto:webpush@example.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "NewVapidKeyError",
			args: args{
				asPrivateKey: "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
				asPublicKey:  "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
				asSubject:    "mailto:webpush@example.com",
			},
			wantErr: assert.Error,
		},
		{
			name: "NewVapidSubjectError",
			args: args{
				asPrivateKey: "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
				asPublicKey:  "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
				asSubject:    "webpush@example.com",
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newVAPID(tt.args.asPrivateKey, tt.args.asPublicKey, tt.args.asSubject)
			if !tt.wantErr(t, err, fmt.Sprintf("newVAPID(%v, %v, %v)", tt.args.asPublicKey, tt.args.asPrivateKey, tt.args.asSubject)) {
				return
			}
		})
	}
}

func Test_vapid_header(t *testing.T) {
	svc, _ := newVAPID(
		"QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
		"BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
		"mailto:webpush@example.com",
	)

	type args struct {
		endpoint string
	}

	tests := []struct {
		name    string
		service *vapid
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "VapidHeaderOK",
			service: svc,
			args: args{
				endpoint: "https://test.com",
			},
			wantErr: assert.NoError,
		},
		{
			name:    "VapidHeaderError",
			service: svc,
			args: args{
				endpoint: "test\ncom",
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.header(tt.args.endpoint)
			if !tt.wantErr(t, err, fmt.Sprintf("header(%v)", tt.args.endpoint)) {
				return
			}
		})
	}
}
