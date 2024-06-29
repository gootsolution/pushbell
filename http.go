package pushbell

import (
	"bytes"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var httpClient = &fasthttp.Client{}

// TODO: make status code mapping

func (s *Service) sendMessage(body *bytes.Buffer, endpoint string, urgency Urgency, ttl time.Duration) (int, error) {
	vapidAuth, err := s.vapid.header(endpoint)
	if err != nil {
		return 0, err
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
		return resp.StatusCode(), err
	}

	return resp.StatusCode(), nil
}
