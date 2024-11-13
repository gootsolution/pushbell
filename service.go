// Package pushbell provides function for sending web push notifications with
// support for the VAPID (Voluntary Application Server Identification)
// specification.

package pushbell

import (
	"fmt"

	"github.com/gootsolution/pushbell/internal/encryption"
	"github.com/gootsolution/pushbell/internal/httpclient"
	"github.com/gootsolution/pushbell/internal/vapid"
)

type Service struct {
	encryption *encryption.Service
	vapid      *vapid.Service
	client     httpclient.Client
	csc        bool
}

// NewService creates new service with given application server keys and subject.
func NewService(options *Options) (*Service, error) {
	vapidService, err := vapid.NewService(
		options.ApplicationServerPublicKey,
		options.ApplicationServerPrivateKey,
		options.ApplicationServerSubject,
	)
	if err != nil {
		return nil, err
	}

	encryptionService, err := encryption.NewService()
	if err != nil {
		return nil, err
	}

	if options.KeyRotationInterval != 0 {
		encryptionService.Rotate(options.KeyRotationInterval)
	}

	client := options.HttpClient
	if client == nil {
		client = httpclient.FastHttp(nil)
	}

	return &Service{
		encryption: encryptionService,
		vapid:      vapidService,
		client:     client,
		csc:        options.CheckStatusCode,
	}, nil
}

// Send sends a WebPush notification with parameters to the specified endpoint.
func (s *Service) Send(push *Push) (int, error) {
	// Cipher text.
	body, err := s.encryption.Encrypt(push.Auth, push.P256DH, push.Plaintext)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt push body: %w", err)
	}

	// Get auth header.
	authHeader, err := s.vapid.Header(push.Endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to generate vapid auth header: %w", err)
	}

	// Prepare headers for client.
	headers := &httpclient.Headers{
		Authorization: authHeader,
		Urgency:       string(push.Urgency),
		TTL:           push.TTL,
	}

	// Request delivery.
	statusCode, err := s.client.RequestDelivery(push.Endpoint, headers, body)
	if err != nil {
		return statusCode, err
	}

	if s.csc {
		return statusCode, CheckStatusCode(statusCode)
	}

	return statusCode, nil
}
