package pushbell

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var httpClient = &fasthttp.Client{}

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

func handleStatusCode(statusCode int) error {
	switch statusCode {
	case fasthttp.StatusOK, fasthttp.StatusCreated, fasthttp.StatusAccepted:
		return nil
	case fasthttp.StatusBadRequest:
		return ErrPushBadRequest
	case fasthttp.StatusUnauthorized:
		return ErrPushUnauthorized
	case fasthttp.StatusForbidden:
		return ErrPushForbidden
	case fasthttp.StatusNotFound:
		return ErrPushNotFound
	case fasthttp.StatusGone:
		return ErrPushGone
	case fasthttp.StatusTooManyRequests:
		return ErrPushTooManyRequests
	case fasthttp.StatusInternalServerError:
		return ErrPushInternalServerError
	case fasthttp.StatusServiceUnavailable:
		return ErrPushServiceUnavailable
	default:
		return ErrPushUnexpectedResponse
	}
}

func (s *Service) sendMessage(endpoint string, urgency Urgency, ttl time.Duration, body *bytes.Buffer) error {
	vapidAuth, err := s.vapid.header(endpoint)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(endpoint)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/octet-stream")
	req.Header.SetContentEncoding("aes128gcm")
	req.Header.SetContentLength(body.Len())
	req.Header.Set("TTL", strconv.FormatInt(int64(ttl/time.Second), 10))
	req.Header.Set("Urgency", string(urgency))
	req.Header.Set("Authorization", vapidAuth)

	req.SetBody(body.Bytes())

	if err = httpClient.Do(req, resp); err != nil {
		return err
	}

	return handleStatusCode(resp.StatusCode())
}
