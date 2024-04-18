package pushbell

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt/v5"
)

type API struct {
	subject      *string
	publicVAPID  *string
	privateVAPID *ecdsa.PrivateKey
}

// New creates a new API client with provided APIConfig.
//
// NOTICE: pushbell use github.com/golang-jwt/jwt/v5 module with
// jwt.MarshalSingleStringAsArray set to true, so be careful, if you use it
// anywhere else in your project.
func New(config *APIConfig) (*API, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	publicVAPID, privateVAPID, err := parseApplicationKeys(config.PublicVAPID, config.PrivateVAPID)
	if err != nil {
		return nil, err
	}

	jwt.MarshalSingleStringAsArray = false

	return &API{
		subject:      config.SubjectVAPID,
		publicVAPID:  publicVAPID,
		privateVAPID: privateVAPID,
	}, nil
}

// Send cipher the plaintext and send it to the endpoint.
// Returns the status code and an error.
func (api *API) Send(endpoint, auth, p256dh string, ttl int, urgency Urgency, plaintext []byte) (int, error) {
	body, err := api.cipherPlaintext(auth, p256dh, plaintext)
	if err != nil {
		return 0, err
	}

	return api.send(endpoint, ttl, urgency, body.Bytes())
}
