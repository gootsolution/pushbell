package pushbell

import "time"

type Service struct {
	vapid      *vapid
	encryption *encryption
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
