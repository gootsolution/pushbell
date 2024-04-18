package pushbell

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

var httpClient = DefaultHTTPClient

// DefaultHTTPClient is configured to be as productive as possible in most scenarios.
var DefaultHTTPClient = &fasthttp.Client{}

// SetHTTPClient allow you to set custom fasthttp.Client. By default,
// DefaultHTTPClient is used. NOTICE: All API instances share the same http
// client, so be careful when having deal with more than one or with already
// running instances.
func SetHTTPClient(client *fasthttp.Client) {
	httpClient = client
}

const headerTemplate = `vapid t=%s, k=%s`

func (api *API) vapidHeader(endpoint string) (string, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Subject:   *api.subject,
		Audience:  jwt.ClaimStrings{fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
	})

	tokenSigned, err := token.SignedString(api.privateVAPID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(headerTemplate, tokenSigned, *api.publicVAPID), nil
}

func (api *API) send(endpoint string, ttl int, urgency Urgency, body []byte) (int, error) {
	authHeader, err := api.vapidHeader(endpoint)
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
	req.Header.SetContentLength(len(body))
	req.Header.Set("TTL", strconv.Itoa(ttl))
	req.Header.Set("Urgency", string(urgency))
	req.Header.Set("Authorization", authHeader)

	req.SetBody(body)

	if err = httpClient.Do(req, resp); err != nil {
		return 0, err
	}

	return resp.StatusCode(), err
}
