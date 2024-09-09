package pushbell

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	applicationServerPrivateKey := "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ"
	applicationServerPublicKey := "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY"
	applicationServerSubject := "mailto:webpush@example.com"

	pb, err := New(applicationServerPrivateKey, applicationServerPublicKey, applicationServerSubject)
	if err != nil {
		panic(err)
	}

	subscriptionEndpoint := "https://fcm.googleapis.com/fcm/send/e2CN0r8ft38:APA91bES3NaBHe_GgsRp_3Ir7f18L38wA5XYRoqZCbjMPEWnkKa07uxheWE5MGZncsPOr0_34zLaFljVqmNqW76KhPSrjdy_pdInnHPEIYAZpdcIYk8oIfo1F_84uKMSqIDXRhngL76S"
	subscriptionAuth := "rm_owGF0xliyVXsrZk1LzQ"
	subscriptionP256DH := "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y"

	message := []byte("{\"title\": \"My first message\"}")

	if err = pb.Send(
		subscriptionEndpoint,
		subscriptionAuth,
		subscriptionP256DH,
		message,
		UrgencyHigh,
		time.Hour,
	); err != nil {
		switch {
		case errors.Is(err, ErrPushGone):
			log.Println(err)
		default:
			panic(err)
		}
	}
}

func TestNew(t *testing.T) {
	type args struct {
		asPrivateKey string
		asPublicKey  string
		asSubject    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "NewServiceOK",
			args: args{
				asPrivateKey: "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
				asPublicKey:  "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
				asSubject:    "mailto:webpush@example.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "NewServiceKeyError",
			args: args{
				asPrivateKey: "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
				asPublicKey:  "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
				asSubject:    "mailto:webpush@example.com",
			},
			wantErr: assert.Error,
		},
		{
			name: "NewServiceSubjectError",
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
			_, err := New(tt.args.asPrivateKey, tt.args.asPublicKey, tt.args.asSubject)
			if !tt.wantErr(t, err, fmt.Sprintf("New(%v, %v, %v)", tt.args.asPublicKey, tt.args.asPrivateKey, tt.args.asSubject)) {
				return
			}
		})
	}
}

func TestService_Send(t *testing.T) {
	ts := setupTestServer()

	svc, _ := New(
		"QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
		"BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
		"mailto:webpush@example.com",
	)

	type args struct {
		endpoint string
		auth     string
		p256dh   string
		message  []byte
		urgency  Urgency
		ttl      time.Duration
	}

	tests := []struct {
		name    string
		service *Service
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "SendOK",
			service: svc,
			args: args{
				endpoint: ts.URL + "/201",
				auth:     "rm_owGF0xliyVXsrZk1LzQ",
				p256dh:   "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y",
				message:  []byte("{\"title\": \"My first message\"}"),
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
			},
			wantErr: assert.NoError,
		},
		{
			name:    "SendBadKey",
			service: svc,
			args: args{
				endpoint: ts.URL + "/201",
				auth:     "rm_owGF0xlsiyVXsrZk1LzQ",
				p256dh:   "BKm5pKbGwksTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y",
				message:  []byte("{\"title\": \"My first message\"}"),
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
			},
			wantErr: assert.Error,
		},
		{
			name:    "SendBadRequest",
			service: svc,
			args: args{
				endpoint: ts.URL + "/400",
				auth:     "rm_owGF0xliyVXsrZk1LzQ",
				p256dh:   "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y",
				message:  []byte("{\"title\": \"My first message\"}"),
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.service.Send(tt.args.endpoint, tt.args.auth, tt.args.p256dh, tt.args.message, tt.args.urgency, tt.args.ttl), fmt.Sprintf("Send(%v, %v, %v, %v, %v, %v)", tt.args.endpoint, tt.args.auth, tt.args.p256dh, tt.args.message, tt.args.urgency, tt.args.ttl))
		})
	}

	ts.Close()
}

func TestService_WithRotation(t *testing.T) {
	contextDone := func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return ctx
	}

	svc, _ := New(
		"QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
		"BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
		"mailto:webpush@example.com",
	)

	type args struct {
		ctx      context.Context
		duration time.Duration
		logger   *slog.Logger
	}

	tests := []struct {
		name    string
		service *Service
		args    args
	}{
		{
			name:    "WithRotationLoggerNil",
			service: svc,
			args: args{
				ctx:      context.Background(),
				duration: time.Second,
				logger:   nil,
			},
		},
		{
			name:    "WithRotationLoggerSlog",
			service: svc,
			args: args{
				ctx:      contextDone(),
				duration: time.Second,
				logger:   slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service.WithRotation(tt.args.ctx, tt.args.duration, tt.args.logger)

			time.Sleep(tt.args.duration * 2)
		})
	}
}
