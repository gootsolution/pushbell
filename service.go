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
	logger     *slog.Logger
}

func NewService(asPublicKey, asPrivateKey, asSubject string) (*Service, error) {
	v, err := newVAPID(asPublicKey, asPrivateKey, asSubject)
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

func (s *Service) Send(endpoint, auth, p256dh string, message []byte, urgency Urgency, ttl time.Duration) (int, error) {
	body, err := s.encryption.encryptMessage(auth, p256dh, message)
	if err != nil {
		return 0, err
	}

	code, err := s.sendMessage(body, endpoint, urgency, ttl)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func (s *Service) WithLogger(logger *slog.Logger) {
	s.logger = logger
}

func (s *Service) WithRotation(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)

	go func(ctx context.Context, ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				if err := s.encryption.rotate(); err != nil {
					if s.logger != nil {
						log.Printf("error while rotating keys: %v", err)
					} else {
						s.logger.Error("error while rotating keys", slog.String("err", err.Error()))
					}
				}
			case <-ctx.Done():
				ticker.Stop()

				return
			}
		}
	}(ctx, ticker)
}
