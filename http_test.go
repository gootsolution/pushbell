package pushbell

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_sendMessage(t *testing.T) {
	ts := setupTestServer()
	svc, err := New(
		"QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ",
		"BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY",
		"mailto:webpush@example.com",
	)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		endpoint string
		urgency  Urgency
		ttl      time.Duration
		body     *bytes.Buffer
	}

	tests := []struct {
		name    string
		service *Service
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "InvalidURL",
			service: svc,
			args: args{
				endpoint: "invalid	url",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.Error,
		},
		{
			name:    "RequestOK",
			service: svc,
			args: args{
				endpoint: ts.URL + "/201",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "RequestUnauthorized",
			service: svc,
			args: args{
				endpoint: ts.URL + "/401",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.Error,
		}, {
			name:    "RequestForbidden",
			service: svc,
			args: args{
				endpoint: ts.URL + "/403",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.Error,
		}, {
			name:    "RequestNotFound",
			service: svc,
			args: args{
				endpoint: ts.URL + "/404",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.Error,
		}, {
			name:    "RequestGone",
			service: svc,
			args: args{
				endpoint: ts.URL + "/410",
				urgency:  UrgencyHigh,
				ttl:      time.Minute,
				body:     new(bytes.Buffer),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.service.sendMessage(tt.args.endpoint, tt.args.urgency, tt.args.ttl, tt.args.body), fmt.Sprintf("sendMessage(%v, %v, %v, %v)", tt.args.endpoint, tt.args.urgency, tt.args.ttl, tt.args.body))
		})
	}

	ts.Close()
}

func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/{code}", func(w http.ResponseWriter, r *http.Request) {
		codeSrt := r.PathValue("code")
		codeInt, _ := strconv.Atoi(codeSrt)

		w.WriteHeader(codeInt)
	})

	return httptest.NewServer(mux)
}
