package pushbell

import (
	"errors"
	"strings"
)

var (
	ErrPublicVAPIDRequired  = errors.New("public VAPID is required")
	ErrPrivateVAPIDRequired = errors.New("private VAPID is required")
	ErrSubjectVAPIDRequired = errors.New("subject VAPID is required")
	ErrSubjectVAPIDFormat   = errors.New("subject VAPID should be either a \"mailto:\" (email) or a \"https:\" URI")
)

type APIConfig struct {
	// PublicVAPID is base64 RawURLEncoded string of ECDH-P256
	PublicVAPID *string

	// PublicVAPID is base64 RawURLEncoded string of ECDH-P256
	PrivateVAPID *string

	// SubjectVAPID should include a contact URI for the application
	// server as either a "mailto:" (email) or a "https:" URI.
	SubjectVAPID *string
}

func (c *APIConfig) validate() error {
	if c.PublicVAPID == nil || *c.PublicVAPID == "" {
		return ErrPublicVAPIDRequired
	}

	if c.PrivateVAPID == nil || *c.PrivateVAPID == "" {
		return ErrPrivateVAPIDRequired
	}

	if c.SubjectVAPID == nil || *c.SubjectVAPID == "" {
		return ErrSubjectVAPIDRequired
	}

	if !strings.HasPrefix(*c.SubjectVAPID, "mailto:") ||
		!strings.HasPrefix(*c.SubjectVAPID, "https:") {
		return ErrSubjectVAPIDFormat
	}

	return nil
}
