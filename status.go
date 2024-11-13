package pushbell

import (
	"errors"
	"net/http"
)

var (
	ErrPushBadRequest          = errors.New("bad request: check the subscription data and message format")
	ErrPushUnauthorized        = errors.New("unauthorized: the request requires authentication")
	ErrPushForbidden           = errors.New("forbidden: the server understood the request, but is refusing to fulfill it")
	ErrPushNotFound            = errors.New("subscription not found: the user may have unsubscribed or the subscription may have expired")
	ErrPushGone                = errors.New("subscription is no longer active: delete it from the server")
	ErrPushTooManyRequests     = errors.New("too many requests: try again later")
	ErrPushInternalServerError = errors.New("internal server error: try again later")
	ErrPushServiceUnavailable  = errors.New("service unavailable: try again later")
	ErrPushUnexpectedResponse  = errors.New("unexpected response from the server")
)

// CheckStatusCode takes an HTTP status code and returns an error if the status code
// indicates an error condition. It can be used to easily handle different types of
// errors that can occur when making HTTP requests.
func CheckStatusCode(statusCode int) error {
	switch statusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil
	case http.StatusBadRequest:
		return ErrPushBadRequest
	case http.StatusUnauthorized:
		return ErrPushUnauthorized
	case http.StatusForbidden:
		return ErrPushForbidden
	case http.StatusNotFound:
		return ErrPushNotFound
	case http.StatusGone:
		return ErrPushGone
	case http.StatusTooManyRequests:
		return ErrPushTooManyRequests
	case http.StatusInternalServerError:
		return ErrPushInternalServerError
	case http.StatusServiceUnavailable:
		return ErrPushServiceUnavailable
	default:
		return ErrPushUnexpectedResponse
	}
}
