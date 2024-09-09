// Package pushbell provides function for sending web push notifications with
// support for the VAPID (Voluntary Application Server Identification)
// specification.

package pushbell

import (
	"context"
	"log"
	"log/slog"
	"time"
)

type Service struct {
	vapid      *vapid
	encryption *encryption
}

// New creates new service with given application server keys and subject.
func New(asPrivateKey, asPublicKey, asSubject string) (*Service, error) {
	v, err := newVAPID(asPrivateKey, asPublicKey, asSubject)
	if err != nil {
		return nil, err
	}

	e, err := newEncryption()
	if err != nil {
		return nil, err
	}

	return &Service{
		encryption: e,
		vapid:      v,
	}, nil
}

func (s *Service) Send(endpoint, auth, p256dh string, message []byte, urgency Urgency, ttl time.Duration) error {
	body, err := s.encryption.encryptMessage(auth, p256dh, message)
	if err != nil {
		return err
	}

	if err = s.sendMessage(endpoint, urgency, ttl, body); err != nil {
		return err
	}

	return nil
}

func (s *Service) WithRotation(ctx context.Context, duration time.Duration, logger *slog.Logger) {
	ticker := time.NewTicker(duration)

	go func(ctx context.Context, ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				if err := s.encryption.rotate(); err != nil {
					if logger != nil {
						log.Printf("error while rotating keys: %v", err)
					} else {
						logger.Error("error while rotating keys", slog.String("err", err.Error()))
					}
				}
			case <-ctx.Done():
				ticker.Stop()

				return
			}
		}
	}(ctx, ticker)
}
